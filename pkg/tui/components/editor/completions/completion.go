package completions

import (
	"github.com/rumpl/rb/pkg/app"
	"github.com/rumpl/rb/pkg/tui/components/completion"
)

type Completion interface {
	Trigger() string
	Items() []completion.Item
	AutoSubmit() bool
	RequiresEmptyEditor() bool
}

func Completions(a *app.App) []Completion {
	return []Completion{
		NewCommandCompletion(a),
		NewFileCompletion(),
	}
}
