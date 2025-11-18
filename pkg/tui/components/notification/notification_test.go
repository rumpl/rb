package notification

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rumpl/rb/pkg/tui/styles"
)

func TestNotification_InitialState(t *testing.T) {
	themeManager := styles.NewManager(styles.ThemeDark)
	n := New(themeManager)

	require.Empty(t, n.items)
	require.False(t, n.Open())
}

func TestNotification_Show(t *testing.T) {
	themeManager := styles.NewManager(styles.ThemeDark)
	n := New(themeManager)

	updated, _ := n.Update(ShowMsg{Text: "Test notification"})

	require.Len(t, updated.items, 1)
	require.Equal(t, "Test notification", updated.items[0].Text)
	require.True(t, updated.Open())
	require.NotEmpty(t, updated.View())
}

func TestNotification_Hide(t *testing.T) {
	themeManager := styles.NewManager(styles.ThemeDark)
	n := New(themeManager)

	updated, _ := n.Update(ShowMsg{Text: "Test"})
	require.Len(t, updated.items, 1)

	updated, _ = updated.Update(HideMsg{})

	require.Empty(t, updated.items)
	require.False(t, updated.Open())
	require.Empty(t, updated.View())
}

func TestNotification_Position(t *testing.T) {
	themeManager := styles.NewManager(styles.ThemeDark)
	n := New(themeManager)
	n.SetSize(100, 50)
	updated, _ := n.Update(ShowMsg{Text: "Test"})
	row, col := updated.position()

	require.Equal(t, 45, row)
	require.Equal(t, 90, col)
}

func TestNotification_GetLayer(t *testing.T) {
	themeManager := styles.NewManager(styles.ThemeDark)
	n := New(themeManager)

	require.Nil(t, n.GetLayer())

	updated, _ := n.Update(ShowMsg{Text: "Test"})
	require.NotNil(t, updated.GetLayer())
}
