package root

import (
	"github.com/spf13/cobra"

	"github.com/rumpl/rb/pkg/agentfile"
	"github.com/rumpl/rb/pkg/config"
	"github.com/rumpl/rb/pkg/mcp"
)

type mcpFlags struct {
	workingDir string
	runConfig  config.RuntimeConfig
}

func newMCPCmd() *cobra.Command {
	var flags mcpFlags

	cmd := &cobra.Command{
		Use:   "mcp <agent-file>|<registry-ref>",
		Short: "Start an agent as an MCP (Model Context Protocol) server",
		Long:  "Start an stdio MCP server that exposes the agent via the Model Context Protocol",
		Example: `  rb mcp ./agent.yaml
  rb mcp ./team.yaml
  rb mcp agentcatalog/pirate`,
		Args:    cobra.ExactArgs(1),
		GroupID: "server",
		RunE:    flags.runMCPCommand,
	}

	cmd.PersistentFlags().StringVar(&flags.workingDir, "working-dir", "", "Set the working directory for the session (applies to tools and relative paths)")
	addRuntimeConfigFlags(cmd, &flags.runConfig)

	return cmd
}

func (f *mcpFlags) runMCPCommand(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	if err := setupWorkingDirectory(f.workingDir); err != nil {
		return err
	}

	agentFilename, err := agentfile.Resolve(ctx, args[0])
	if err != nil {
		return err
	}

	return mcp.StartMCPServer(ctx, agentFilename, f.runConfig)
}
