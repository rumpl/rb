package root

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/rumpl/rb/pkg/version"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Long:  "Display the version and commit hash",
		Args:  cobra.NoArgs,
		Run:   runVersionCommand,
	}
}

func runVersionCommand(cmd *cobra.Command, args []string) {
	fmt.Printf("rb version %s\n", version.Version)
	fmt.Printf("Commit: %s\n", version.Commit)
}
