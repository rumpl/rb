package root

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/rumpl/rb/pkg/content"
	"github.com/rumpl/rb/pkg/oci"
	"github.com/rumpl/rb/pkg/remote"
)

func newPushCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "push <agent-file> <registry-ref>",
		Short:   "Push an agent to an OCI registry",
		Long:    "Push an agent configuration file to an OCI registry",
		GroupID: "core",
		Args:    cobra.ExactArgs(2),
		RunE:    runPushCommand,
	}
}

func runPushCommand(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	tag := args[1]

	store, err := content.NewStore()
	if err != nil {
		return err
	}

	_, err = oci.PackageFileAsOCIToStore(filePath, tag, store)
	if err != nil {
		return fmt.Errorf("failed to build artifact: %w", err)
	}

	slog.Debug("Starting push", "registry_ref", tag)

	fmt.Printf("Pushing agent %s to %s\n", filePath, tag)

	err = remote.Push(tag)
	if err != nil {
		return fmt.Errorf("failed to push artifact: %w", err)
	}

	fmt.Printf("Successfully pushed artifact to %s\n", tag)
	return nil
}
