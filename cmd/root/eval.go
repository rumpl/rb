package root

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/rumpl/rb/pkg/agentfile"
	"github.com/rumpl/rb/pkg/config"
	"github.com/rumpl/rb/pkg/evaluation"
	"github.com/rumpl/rb/pkg/teamloader"
)

type evalFlags struct {
	runConfig config.RuntimeConfig
}

func newEvalCmd() *cobra.Command {
	var flags evalFlags

	cmd := &cobra.Command{
		Use:     "eval <agent-file>|<registry-ref> <eval-dir>",
		Short:   "Run evaluations for an agent",
		GroupID: "advanced",
		Args:    cobra.ExactArgs(2),
		RunE:    flags.runEvalCommand,
	}

	addRuntimeConfigFlags(cmd, &flags.runConfig)

	return cmd
}

func (f *evalFlags) runEvalCommand(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	agentFilename, err := agentfile.Resolve(ctx, args[0])
	if err != nil {
		return err
	}

	agents, err := teamloader.Load(cmd.Context(), agentFilename, f.runConfig)
	if err != nil {
		return err
	}

	evalResults, err := evaluation.Evaluate(cmd.Context(), agents, args[1])
	if err != nil {
		return err
	}

	for _, evalResult := range evalResults {
		fmt.Printf("Eval file: %s\n", evalResult.EvalFile)
		fmt.Printf("Tool trajectory score: %f\n", evalResult.Score.ToolTrajectoryScore)
		fmt.Printf("Rouge-1 score: %f\n", evalResult.Score.Rouge1Score)
	}

	return nil
}
