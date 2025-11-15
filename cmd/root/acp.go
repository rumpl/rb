package root

import (
	"log/slog"

	acpsdk "github.com/coder/acp-go-sdk"
	"github.com/spf13/cobra"

	"github.com/rumpl/rb/pkg/acp"
	"github.com/rumpl/rb/pkg/agentfile"
	"github.com/rumpl/rb/pkg/config"
)

type acpFlags struct {
	runConfig config.RuntimeConfig
}

func newACPCmd() *cobra.Command {
	var flags acpFlags

	cmd := &cobra.Command{
		Use:   "acp <agent-file>|<registry-ref>",
		Short: "Start an agent as an ACP (Agent Client Protocol) server",
		Long:  "Start an ACP server that exposes the agent via the Agent Client Protocol",
		Example: `  rb acp ./agent.yaml
  rb acp ./team.yaml
  rb acp agentcatalog/pirate`,
		Args:    cobra.ExactArgs(1),
		GroupID: "server",
		RunE:    flags.runACPCommand,
	}

	addRuntimeConfigFlags(cmd, &flags.runConfig)

	return cmd
}

func (f *acpFlags) runACPCommand(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	agentFilename, err := agentfile.Resolve(ctx, args[0])
	if err != nil {
		return err
	}

	slog.Debug("Starting ACP server", "agent_file", agentFilename)

	acpAgent := acp.NewAgent(agentFilename, f.runConfig)
	conn := acpsdk.NewAgentSideConnection(acpAgent, cmd.OutOrStdout(), cmd.InOrStdin())
	conn.SetLogger(slog.Default())
	acpAgent.SetAgentConnection(conn)
	defer acpAgent.Stop(ctx)

	slog.Debug("acp started, waiting for conn")
	<-conn.Done()

	return nil
}
