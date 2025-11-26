package ghosttyhelp

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/kregenrek/tmuxman/internal/tui/theme"
)

type shortcut struct {
	key  string
	desc string
}

// Model renders a list of Ghostty -> tmux shortcuts.
type Model struct {
	width  int
	height int
}

var shortcuts = []shortcut{
	{"Cmd+H/J/K/L", "Navigate panes"},
	{"Cmd+[ / ]", "Prev/next window"},
	{"Cmd+T", "New window"},
	{"Cmd+W", "Close window"},
	{"Cmd+1…9", "Jump to window"},
	{"Cmd+R", "Respawn pane"},
	{"Cmd+Shift+W", "Kill session"},
	{"Cmd+Shift+H/J/K/L", "Resize panes"},
	{"Cmd+Backspace", "Clear line"},
	{"Cmd+Shift+P", "Command palette"},
	{"Cmd+I", "Toggle this help"},
}

// NewModel creates a help view with the predefined shortcuts.
func NewModel() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd { return tea.ClearScreen }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m Model) View() string {
	var b strings.Builder

	// Title - using centralized theme
	b.WriteString(theme.HelpTitle.Render("⌨️  Ghostty → tmux"))
	b.WriteString("\n\n")

	// Shortcuts - using centralized theme
	for _, s := range shortcuts {
		b.WriteString(theme.ShortcutKey.Render(s.key))
		b.WriteString(theme.ShortcutDesc.Render(s.desc))
		b.WriteString("\n")
	}

	// Footer note
	b.WriteString("\n")
	b.WriteString(theme.ShortcutNote.Render("Cmd sends tmux prefix automatically"))
	b.WriteString("\n\n")

	// Close hint
	b.WriteString(theme.ShortcutHint.Render("esc to close"))

	return b.String()
}
