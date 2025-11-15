package root

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"

	"github.com/rumpl/rb/pkg/agentfile"
	"github.com/rumpl/rb/pkg/app"
	"github.com/rumpl/rb/pkg/config"
	"github.com/rumpl/rb/pkg/runtime"
	"github.com/rumpl/rb/pkg/session"
	"github.com/rumpl/rb/pkg/team"
	"github.com/rumpl/rb/pkg/teamloader"
	"github.com/rumpl/rb/pkg/tui"
)

type runExecFlags struct {
	agentName      string
	workingDir     string
	autoApprove    bool
	attachmentPath string
	remoteAddress  string
	modelOverrides []string
	runConfig      config.RuntimeConfig
}

func newRunCmd() *cobra.Command {
	var flags runExecFlags

	cmd := &cobra.Command{
		Use:   "run <agent-file>|<registry-ref> [message|-]",
		Short: "Run an agent",
		Long:  "Run an agent with the specified configuration and prompt",
		Example: `  rb run ./agent.yaml
  rb run ./team.yaml --agent root
  rb run ./echo.yaml "INSTRUCTIONS"
  echo "INSTRUCTIONS" | rb run ./echo.yaml -`,
		GroupID: "core",
		Args:    cobra.RangeArgs(1, 2),
		RunE:    flags.runRunCommand,
	}

	addRunOrExecFlags(cmd, &flags)
	addRuntimeConfigFlags(cmd, &flags.runConfig)

	return cmd
}

func addRunOrExecFlags(cmd *cobra.Command, flags *runExecFlags) {
	cmd.PersistentFlags().StringVarP(&flags.agentName, "agent", "a", "root", "Name of the agent to run")
	cmd.PersistentFlags().StringVar(&flags.workingDir, "working-dir", "", "Set the working directory for the session (applies to tools and relative paths)")
	cmd.PersistentFlags().BoolVar(&flags.autoApprove, "yolo", false, "Automatically approve all tool calls without prompting")
	cmd.PersistentFlags().StringVar(&flags.attachmentPath, "attach", "", "Attach an image file to the message")
	cmd.PersistentFlags().StringArrayVar(&flags.modelOverrides, "model", nil, "Override agent model: [agent=]provider/model (repeatable)")
	cmd.PersistentFlags().StringVar(&flags.remoteAddress, "remote", "", "Use remote runtime with specified address")
}

func (f *runExecFlags) runRunCommand(cmd *cobra.Command, args []string) error {
	return f.run(cmd.Context(), args)
}

func (f *runExecFlags) run(ctx context.Context, args []string) error {
	slog.Debug("Starting agent", "agent", f.agentName)

	if err := f.setupWorkingDirectory(); err != nil {
		return err
	}

	agentFileName := ""

	var rt runtime.Runtime
	var sess *session.Session
	var err error
	switch {
	case f.remoteAddress != "":
		agentFileName = args[0]
		rt, sess, err = f.createRemoteRuntimeAndSession(ctx, agentFileName)
		if err != nil {
			return err
		}

	default:
		agentFileName, err = f.resolveAgentFile(ctx, args[0])
		if err != nil {
			return err
		}

		t, err := f.loadAgentFrom(ctx, teamloader.NewFileSource(agentFileName))
		if err != nil {
			return err
		}

		rt, sess, err = f.createLocalRuntimeAndSession(t)
		if err != nil {
			return err
		}
	}

	return handleRunMode(ctx, agentFileName, rt, sess, args)
}

func (f *runExecFlags) setupWorkingDirectory() error {
	return setupWorkingDirectory(f.workingDir)
}

// resolveAgentFile is a wrapper method that calls the agentfile.Resolve function
// after checking for remote address
func (f *runExecFlags) resolveAgentFile(ctx context.Context, agentFilename string) (string, error) {
	if f.remoteAddress != "" {
		return agentFilename, nil
	}
	return agentfile.Resolve(ctx, agentFilename)
}

func (f *runExecFlags) loadAgentFrom(ctx context.Context, source teamloader.AgentSource) (*team.Team, error) {
	t, err := teamloader.LoadFrom(ctx, source, f.runConfig, teamloader.WithModelOverrides(f.modelOverrides))
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		if err := t.StopToolSets(ctx); err != nil {
			slog.Error("Failed to stop tool sets", "error", err)
		}
	}()

	return t, nil
}

func (f *runExecFlags) createRemoteRuntimeAndSession(ctx context.Context, originalFilename string) (runtime.Runtime, *session.Session, error) {
	remoteClient, err := runtime.NewClient(f.remoteAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create remote client: %w", err)
	}

	sessTemplate := session.New(
		session.WithToolsApproved(f.autoApprove),
	)

	sess, err := remoteClient.CreateSession(ctx, sessTemplate)
	if err != nil {
		return nil, nil, err
	}

	remoteRt, err := runtime.NewRemoteRuntime(remoteClient,
		runtime.WithRemoteCurrentAgent(f.agentName),
		runtime.WithRemoteAgentFilename(originalFilename),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create remote runtime: %w", err)
	}

	slog.Debug("Using remote runtime", "address", f.remoteAddress, "agent", f.agentName)
	return remoteRt, sess, nil
}

func (f *runExecFlags) createLocalRuntimeAndSession(t *team.Team) (runtime.Runtime, *session.Session, error) {
	agent, err := t.Agent(f.agentName)
	if err != nil {
		return nil, nil, err
	}

	sess := session.New(
		session.WithMaxIterations(agent.MaxIterations()),
		session.WithToolsApproved(f.autoApprove),
	)

	localRt, err := runtime.New(t,
		runtime.WithCurrentAgent(f.agentName),
		runtime.WithTracer(otel.Tracer(AppName)),
		runtime.WithRootSessionID(sess.ID),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create runtime: %w", err)
	}

	slog.Debug("Using local runtime", "agent", f.agentName)
	return localRt, sess, nil
}

func readInitialMessage(args []string) (*string, error) {
	if len(args) < 2 {
		return nil, nil
	}

	if args[1] == "-" {
		buf, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("failed to read from stdin: %w", err)
		}
		text := string(buf)
		return &text, nil
	}

	return &args[1], nil
}

func handleRunMode(ctx context.Context, agentFilename string, rt runtime.Runtime, sess *session.Session, args []string) error {
	firstMessage, err := readInitialMessage(args)
	if err != nil {
		return err
	}

	a := app.New(agentFilename, rt, sess, firstMessage)
	m := tui.New(a)

	progOpts := []tea.ProgramOption{
		tea.WithContext(ctx),
		tea.WithFilter(tui.MouseEventFilter),
	}

	p := tea.NewProgram(m, progOpts...)

	go a.Subscribe(ctx, p)

	_, err = p.Run()
	return err
}
