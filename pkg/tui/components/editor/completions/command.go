package completions

import (
	"context"

	"github.com/rumpl/rb/pkg/app"
	"github.com/rumpl/rb/pkg/tui/commands"
	"github.com/rumpl/rb/pkg/tui/components/completion"
)

type commandCompletion struct {
	app *app.App
}

func NewCommandCompletion(a *app.App) Completion {
	return &commandCompletion{
		app: a,
	}
}

func (c *commandCompletion) AutoSubmit() bool {
	return true
}

func (c *commandCompletion) RequiresEmptyEditor() bool {
	return true
}

func (c *commandCompletion) Trigger() string {
	return "/"
}

func (c *commandCompletion) Items() []completion.Item {
	cmds := commands.BuildCommandCategories(context.Background(), c.app)
	items := make([]completion.Item, 0, len(cmds))
	for _, cmd := range cmds {
		for _, command := range cmd.Commands {
			items = append(items, completion.Item{
				Label:       command.Label,
				Description: command.Description,
				Value:       command.SlashCommand,
				Execute:     command.Execute,
			})
		}
	}
	return items
}
