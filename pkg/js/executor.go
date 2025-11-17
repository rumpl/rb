package js

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/dop251/goja"

	"github.com/rumpl/rb/pkg/tools"
)

// createToolWrapper creates a JavaScript-callable wrapper for a Go tool.
// This wrapper handles argument marshaling, optional parameter filtering,
// and tool execution.
func createToolWrapper(ctx context.Context, tool tools.Tool) func(args map[string]any) (string, error) {
	return func(args map[string]any) (string, error) {
		// Extract required parameters from the tool schema
		var toolArgs struct {
			Required []string `json:"required"`
		}

		if err := tools.ConvertSchema(tool.Parameters, &toolArgs); err != nil {
			return "", fmt.Errorf("parsing tool schema: %w", err)
		}

		// Filter out nil values for optional parameters
		nonNilArgs := make(map[string]any)
		for k, v := range args {
			if slices.Contains(toolArgs.Required, k) || v != nil {
				nonNilArgs[k] = v
			}
		}

		// Marshal arguments to JSON
		arguments, err := json.Marshal(nonNilArgs)
		if err != nil {
			return "", fmt.Errorf("marshaling tool arguments: %w", err)
		}

		// Call the tool handler
		result, err := tool.Handler(ctx, tools.ToolCall{
			Function: tools.FunctionCall{
				Name:      tool.Name,
				Arguments: string(arguments),
			},
		})
		if err != nil {
			return "", fmt.Errorf("executing tool %s: %w", tool.Name, err)
		}

		return result.Output, nil
	}
}

// InjectTools injects all tools from the given toolsets into the JavaScript VM.
// Each tool becomes a callable JavaScript function.
func InjectTools(ctx context.Context, vm *goja.Runtime, toolSets []tools.ToolSet) error {
	for _, toolSet := range toolSets {
		allTools, err := toolSet.Tools(ctx)
		if err != nil {
			return fmt.Errorf("getting tools from toolset: %w", err)
		}

		for _, tool := range allTools {
			if err := vm.Set(tool.Name, createToolWrapper(ctx, tool)); err != nil {
				return fmt.Errorf("injecting tool %s: %w", tool.Name, err)
			}
		}
	}

	return nil
}
