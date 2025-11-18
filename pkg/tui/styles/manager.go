package styles

import (
	"sync"

	"github.com/alecthomas/chroma/v2"
)

// Manager manages the current theme and provides thread-safe access to styles.
type Manager struct {
	mu           sync.RWMutex
	currentTheme ThemeName
	theme        Theme
}

// NewManager creates a new theme manager with the specified theme.
func NewManager(themeName ThemeName) *Manager {
	return &Manager{
		currentTheme: themeName,
		theme:        GetTheme(themeName),
	}
}

// GetTheme returns the current theme.
func (m *Manager) GetTheme() Theme {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.theme
}

// SetTheme sets the current theme.
func (m *Manager) SetTheme(themeName ThemeName) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.currentTheme = themeName
	m.theme = GetTheme(themeName)
}

// CurrentThemeName returns the name of the current theme.
func (m *Manager) CurrentThemeName() ThemeName {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentTheme
}

// ChromaStyle returns the Chroma style for the current theme.
func (m *Manager) ChromaStyle() *chroma.Style {
	return m.GetTheme().GetChromaStyle()
}
