package agent

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rumpl/rb/pkg/tools"
)

// mockToolSet is a test toolset that provides test tools
type mockToolSet struct {
	tools []tools.Tool
}

func (m *mockToolSet) Tools(context.Context) ([]tools.Tool, error) {
	return m.tools, nil
}

func (m *mockToolSet) Instructions() string {
	return ""
}

func (m *mockToolSet) Start(context.Context) error {
	return nil
}

func (m *mockToolSet) Stop(context.Context) error {
	return nil
}

func (m *mockToolSet) SetElicitationHandler(tools.ElicitationHandler) {}

func (m *mockToolSet) SetOAuthSuccessHandler(func()) {}

func TestExpandInstruction_NoTemplateExpression(t *testing.T) {
	instruction := "This is a simple instruction without any templates"

	expanded, err := ExpandInstruction(t.Context(), instruction, nil)
	require.NoError(t, err)
	assert.Equal(t, instruction, expanded, "instruction without templates should remain unchanged")
}

func TestExpandInstruction_SimpleTool(t *testing.T) {
	// Create a simple test tool that returns a fixed value
	testTool := tools.Tool{
		Name:        "test_tool",
		Description: "A test tool",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"message": map[string]any{
					"type": "string",
				},
			},
			"required": []any{"message"},
		},
		Handler: func(ctx context.Context, toolCall tools.ToolCall) (*tools.ToolCallResult, error) {
			var args struct {
				Message string `json:"message"`
			}
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
				return nil, err
			}
			return &tools.ToolCallResult{
				Output: "Tool executed: " + args.Message,
			}, nil
		},
	}

	toolSet := &mockToolSet{tools: []tools.Tool{testTool}}

	instruction := `You are a helpful assistant. ${test_tool({message: "hello"})}`

	expanded, err := ExpandInstruction(t.Context(), instruction, []tools.ToolSet{toolSet})
	require.NoError(t, err)
	assert.Equal(t, "You are a helpful assistant. Tool executed: hello", expanded)
}

func TestExpandInstruction_MultipleTool(t *testing.T) {
	// Create tools
	tool1 := tools.Tool{
		Name:        "get_greeting",
		Description: "Returns a greeting",
		Parameters: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		},
		Handler: func(ctx context.Context, toolCall tools.ToolCall) (*tools.ToolCallResult, error) {
			return &tools.ToolCallResult{Output: "Hello"}, nil
		},
	}

	tool2 := tools.Tool{
		Name:        "get_name",
		Description: "Returns a name",
		Parameters: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		},
		Handler: func(ctx context.Context, toolCall tools.ToolCall) (*tools.ToolCallResult, error) {
			return &tools.ToolCallResult{Output: "World"}, nil
		},
	}

	toolSet := &mockToolSet{tools: []tools.Tool{tool1, tool2}}

	instruction := `${get_greeting()}, ${get_name()}!`

	expanded, err := ExpandInstruction(t.Context(), instruction, []tools.ToolSet{toolSet})
	require.NoError(t, err)
	assert.Equal(t, "Hello, World!", expanded)
}

func TestExpandInstruction_ToolWithOptionalParameters(t *testing.T) {
	ctx := t.Context()

	testTool := tools.Tool{
		Name:        "format_message",
		Description: "Formats a message",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"text": map[string]any{
					"type": "string",
				},
				"prefix": map[string]any{
					"type": "string",
				},
			},
			"required": []any{"text"},
		},
		Handler: func(ctx context.Context, toolCall tools.ToolCall) (*tools.ToolCallResult, error) {
			var args struct {
				Text   string `json:"text"`
				Prefix string `json:"prefix,omitempty"`
			}
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
				return nil, err
			}
			result := args.Text
			if args.Prefix != "" {
				result = args.Prefix + ": " + result
			}
			return &tools.ToolCallResult{Output: result}, nil
		},
	}

	toolSet := &mockToolSet{tools: []tools.Tool{testTool}}

	// Test with optional parameter
	instruction1 := `${format_message({text: "test", prefix: "INFO"})}`
	expanded1, err := ExpandInstruction(ctx, instruction1, []tools.ToolSet{toolSet})
	require.NoError(t, err)
	assert.Equal(t, "INFO: test", expanded1)

	// Test without optional parameter
	instruction2 := `${format_message({text: "test"})}`
	expanded2, err := ExpandInstruction(ctx, instruction2, []tools.ToolSet{toolSet})
	require.NoError(t, err)
	assert.Equal(t, "test", expanded2)
}

func TestExpandInstruction_InvalidSyntax(t *testing.T) {
	testTool := tools.Tool{
		Name:        "test_tool",
		Description: "A test tool",
		Parameters: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		},
		Handler: func(ctx context.Context, toolCall tools.ToolCall) (*tools.ToolCallResult, error) {
			return &tools.ToolCallResult{Output: "result"}, nil
		},
	}

	toolSet := &mockToolSet{tools: []tools.Tool{testTool}}

	// Invalid JavaScript syntax
	instruction := `${invalid javascript here}`

	_, err := ExpandInstruction(t.Context(), instruction, []tools.ToolSet{toolSet})
	assert.Error(t, err, "invalid syntax should return an error")
}

func TestExpandInstruction_ToolNotFound(t *testing.T) {
	toolSet := &mockToolSet{tools: []tools.Tool{}}

	instruction := `${non_existent_tool()}`

	_, err := ExpandInstruction(t.Context(), instruction, []tools.ToolSet{toolSet})
	assert.Error(t, err, "calling non-existent tool should return an error")
}

func TestExpandInstruction_ComplexExpression(t *testing.T) {
	getTool := tools.Tool{
		Name:        "get_value",
		Description: "Gets a value",
		Parameters: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		},
		Handler: func(ctx context.Context, toolCall tools.ToolCall) (*tools.ToolCallResult, error) {
			return &tools.ToolCallResult{Output: "42"}, nil
		},
	}

	toolSet := &mockToolSet{tools: []tools.Tool{getTool}}

	// Test with expression that includes JavaScript operations
	instruction := `The answer is ${get_value()} and double is ${parseInt(get_value()) * 2}`

	expanded, err := ExpandInstruction(t.Context(), instruction, []tools.ToolSet{toolSet})
	require.NoError(t, err)
	assert.Equal(t, "The answer is 42 and double is 84", expanded)
}

func TestExpandInstruction_EmptyToolSet(t *testing.T) {
	instruction := "No tools needed here"

	expanded, err := ExpandInstruction(t.Context(), instruction, []tools.ToolSet{})
	require.NoError(t, err)
	assert.Equal(t, instruction, expanded)
}

func TestAgent_InstructionExpansion(t *testing.T) {
	ctx := t.Context()

	testTool := tools.Tool{
		Name:        "get_config",
		Description: "Gets configuration",
		Parameters: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		},
		Handler: func(ctx context.Context, toolCall tools.ToolCall) (*tools.ToolCallResult, error) {
			return &tools.ToolCallResult{Output: "config_value"}, nil
		},
	}

	toolSet := &mockToolSet{tools: []tools.Tool{testTool}}

	agent := New("test", "Use this config: ${get_config()}",
		WithToolSets(toolSet))

	// First call should expand the instruction
	instruction1 := agent.Instruction(ctx)
	assert.Equal(t, "Use this config: config_value", instruction1)

	// Second call should return cached result
	instruction2 := agent.Instruction(ctx)
	assert.Equal(t, instruction1, instruction2, "instruction should be cached")
}

func TestAgent_InstructionExpansion_NoTemplate(t *testing.T) {
	agent := New("test", "Simple instruction without templates")

	instruction := agent.Instruction(t.Context())
	assert.Equal(t, "Simple instruction without templates", instruction)
}

func TestAgent_InstructionExpansion_Error(t *testing.T) {
	// Agent with invalid template expression
	agent := New("test", "Invalid: ${non_existent_tool()}")

	// Should return original instruction on error
	instruction := agent.Instruction(t.Context())
	assert.Equal(t, "Invalid: ${non_existent_tool()}", instruction)

	// Check that warnings were generated
	warnings := agent.DrainWarnings()
	assert.NotEmpty(t, warnings, "should have warnings about expansion failure")
}

func TestExpandInstruction_WithJSONResult(t *testing.T) {
	jsonTool := tools.Tool{
		Name:        "get_json",
		Description: "Returns JSON data",
		Parameters: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		},
		Handler: func(ctx context.Context, toolCall tools.ToolCall) (*tools.ToolCallResult, error) {
			return &tools.ToolCallResult{
				Output: `{"key": "value", "number": 123}`,
			}, nil
		},
	}

	toolSet := &mockToolSet{tools: []tools.Tool{jsonTool}}

	instruction := `Here is the data: ${get_json()}`

	expanded, err := ExpandInstruction(t.Context(), instruction, []tools.ToolSet{toolSet})
	require.NoError(t, err)
	assert.Contains(t, expanded, `{"key": "value", "number": 123}`)
}

func TestExpandInstruction_MultilineResult(t *testing.T) {
	multilineTool := tools.Tool{
		Name:        "get_multiline",
		Description: "Returns multiline text",
		Parameters: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		},
		Handler: func(ctx context.Context, toolCall tools.ToolCall) (*tools.ToolCallResult, error) {
			return &tools.ToolCallResult{
				Output: "line1\nline2\nline3",
			}, nil
		},
	}

	toolSet := &mockToolSet{tools: []tools.Tool{multilineTool}}

	instruction := `Files:\n${get_multiline()}`

	expanded, err := ExpandInstruction(t.Context(), instruction, []tools.ToolSet{toolSet})
	require.NoError(t, err)
	assert.Equal(t, "Files:\nline1\nline2\nline3", expanded)
}
