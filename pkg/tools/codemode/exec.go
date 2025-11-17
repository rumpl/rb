package codemode

import (
	"bytes"
	"context"
	"fmt"

	"github.com/dop251/goja"

	"github.com/rumpl/rb/pkg/js"
)

type ScriptResult struct {
	Value  string `json:"value" jsonschema:"The value returned by the script"`
	StdOut string `json:"stdout" jsonschema:"The standard output of the console"`
	StdErr string `json:"stderr" jsonschema:"The standard error of the console"`
}

func (c *codeModeTool) runJavascript(ctx context.Context, script string) (ScriptResult, error) {
	vm := goja.New()

	// Inject console object to the help the LLM debug its own code.
	var (
		stdOut bytes.Buffer
		stdErr bytes.Buffer
	)
	_ = vm.Set("console", console(&stdOut, &stdErr))

	// Inject every tool as a javascript function.
	if err := js.InjectTools(ctx, vm, c.toolsets); err != nil {
		return ScriptResult{}, err
	}

	// Wrap the user script in an IIFE to allow top-level returns.
	script = "(() => {\n" + script + "\n})()"

	// Run the script.
	v, err := vm.RunString(script)
	if err != nil {
		return ScriptResult{
			StdOut: stdOut.String(),
			StdErr: stdErr.String(),
			Value:  err.Error(),
		}, nil
	}

	value := ""
	if result := v.Export(); result != nil {
		value = fmt.Sprintf("%v", result)
	}

	return ScriptResult{
		StdOut: stdOut.String(),
		StdErr: stdErr.String(),
		Value:  value,
	}, nil
}
