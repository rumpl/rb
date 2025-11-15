package root

import (
	"github.com/spf13/cobra"
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
	// TODO: ???
	return f.run(cmd.Context(), args)
}
