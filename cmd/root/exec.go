package root

import (
	"github.com/spf13/cobra"

	"github.com/rumpl/rb/pkg/cli"
	"github.com/rumpl/rb/pkg/telemetry"
)

func newExecCmd() *cobra.Command {
	var flags runExecFlags

	cmd := &cobra.Command{
		Use:   "exec <agent-file>|<registry-ref>",
		Short: "Execute an agent",
		Long:  "Execute an agent (Single user message / No TUI)",
		Example: `  rb exec ./agent.yaml
  rb exec ./team.yaml --agent root
  rb exec ./echo.yaml "INSTRUCTIONS"
  echo "INSTRUCTIONS" | rb exec ./echo.yaml -`,
		GroupID: "core",
		Args:    cobra.RangeArgs(1, 2),
		RunE:    flags.runExecCommand,
	}

	addRunOrExecFlags(cmd, &flags)
	addRuntimeConfigFlags(cmd, &flags.runConfig)

	return cmd
}

func (f *runExecFlags) runExecCommand(cmd *cobra.Command, args []string) error {
	telemetry.TrackCommand("exec", args)

	ctx := cmd.Context()
	out := cli.NewPrinter(cmd.OutOrStdout())

	return f.runOrExec(ctx, out, args, true)
}
