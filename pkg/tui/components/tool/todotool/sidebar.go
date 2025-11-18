package todotool

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rumpl/rb/pkg/tools"
	"github.com/rumpl/rb/pkg/tools/builtin"
	"github.com/rumpl/rb/pkg/tui/service"
	"github.com/rumpl/rb/pkg/tui/styles"
	"github.com/rumpl/rb/pkg/tui/types"
)

// SidebarComponent represents the todo display component for the sidebar
type SidebarComponent struct {
	manager      *service.TodoManager
	width        int
	themeManager *styles.Manager
}

func NewSidebarComponent(manager *service.TodoManager, themeManager *styles.Manager) *SidebarComponent {
	return &SidebarComponent{
		manager:      manager,
		width:        20,
		themeManager: themeManager,
	}
}

func (c *SidebarComponent) SetSize(width int) {
	c.width = width
}

func (c *SidebarComponent) SetTodos(toolCall tools.ToolCall) error {
	params, err := parseTodoArgs(toolCall)
	if err != nil {
		return err
	}

	toolName := toolCall.Function.Name
	switch toolName {
	case builtin.ToolNameCreateTodo:
		p := params.(builtin.CreateTodoArgs)
		newID := generateTodoID(c.manager.GetTodos())
		c.manager.AddTodo(newID, p.Description, "pending")

	case builtin.ToolNameCreateTodos:
		p := params.(builtin.CreateTodosArgs)
		for _, desc := range p.Descriptions {
			newID := generateTodoID(c.manager.GetTodos())
			c.manager.AddTodo(newID, desc, "pending")
		}

	case builtin.ToolNameUpdateTodo:
		p := params.(builtin.UpdateTodoArgs)
		c.manager.UpdateTodo(p.ID, p.Status)
	}

	return nil
}

func (c *SidebarComponent) Render() string {
	theme := c.themeManager.GetTheme()
	if len(c.manager.GetTodos()) == 0 {
		return ""
	}

	var content strings.Builder
	content.WriteString(theme.HighlightStyle.Render("TODOs"))
	content.WriteString("\n")

	for _, todo := range c.manager.GetTodos() {
		content.WriteString(c.renderTodoLine(todo, c.width))
		content.WriteString("\n")
	}

	return content.String()
}

func (c *SidebarComponent) renderTodoLine(todo types.Todo, maxWidth int) string {
	theme := c.themeManager.GetTheme()
	icon, style := renderTodoIcon(todo.Status, &theme)

	description := todo.Description
	maxDescWidth := max(maxWidth-2, 3)
	if len(description) > maxDescWidth {
		description = description[:maxDescWidth-3] + "..."
	}

	styledIcon := style.Render(icon)
	styledDescription := style.Render(description)
	return fmt.Sprintf("%s %s", styledIcon, styledDescription)
}

func parseTodoArgs(toolCall tools.ToolCall) (any, error) {
	toolName := toolCall.Function.Name
	arguments := toolCall.Function.Arguments

	switch toolName {
	case builtin.ToolNameCreateTodo:
		var params builtin.CreateTodoArgs
		if err := json.Unmarshal([]byte(arguments), &params); err != nil {
			return nil, err
		}
		return params, nil
	case builtin.ToolNameCreateTodos:
		var params builtin.CreateTodosArgs
		if err := json.Unmarshal([]byte(arguments), &params); err != nil {
			return nil, err
		}
		return params, nil
	case builtin.ToolNameUpdateTodo:
		var params builtin.UpdateTodoArgs
		if err := json.Unmarshal([]byte(arguments), &params); err != nil {
			return nil, err
		}
		return params, nil
	case builtin.ToolNameListTodos:
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown tool name: %s", toolName)
	}
}

func generateTodoID(todos []types.Todo) string {
	return fmt.Sprintf("todo_%d", len(todos)+1)
}
