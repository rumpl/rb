package agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/dop251/goja"

	"github.com/rumpl/rb/pkg/js"
	"github.com/rumpl/rb/pkg/tools"
)

// ExpandInstruction expands JavaScript template literal expressions in the instruction string.
// It supports ${expression} syntax where expression can call any available tool.
// Example: "You have access to these files: ${list_directory({path: '.'})}"
func ExpandInstruction(ctx context.Context, instruction string, toolSets []tools.ToolSet) (string, error) {
	// Fast path: if no template expressions, return as-is
	if !strings.Contains(instruction, "${") {
		return instruction, nil
	}

	vm := goja.New()

	// Inject all tools as JavaScript functions
	if err := js.InjectTools(ctx, vm, toolSets); err != nil {
		return "", err
	}

	// Wrap in backticks to create a template literal
	script := "`" + instruction + "`"

	// Execute the template literal
	result, err := vm.RunString(script)
	if err != nil {
		return "", fmt.Errorf("executing instruction template: %w", err)
	}

	// Convert result to string
	expanded := fmt.Sprintf("%v", result.Export())

	return expanded, nil
}
