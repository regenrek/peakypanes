package peakypanes

import (
	"testing"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// TestStatusIcon tests the status icon helper function
func TestStatusIcon(t *testing.T) {
	tests := []struct {
		name   string
		status Status
		want   string
	}{
		{name: "current", status: StatusCurrent, want: "◆"},
		{name: "running", status: StatusRunning, want: "●"},
		{name: "stopped", status: StatusStopped, want: "○"},
		{name: "unknown", status: Status(99), want: "?"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := statusIcon(tt.status)
			if got != tt.want {
				t.Errorf("statusIcon(%d) = %q, want %q", tt.status, got, tt.want)
			}
		})
	}
}

// TestExpandPath tests the path expansion helper
func TestExpandPath(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string // Empty means check for home prefix
	}{
		{name: "empty", input: "", want: ""},
		{name: "absolute", input: "/tmp/test", want: "/tmp/test"},
		{name: "relative", input: "test/path", want: "test/path"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandPath(tt.input)
			if tt.want != "" && got != tt.want {
				t.Errorf("expandPath(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestExpandPathTilde tests tilde expansion specifically
func TestExpandPathTilde(t *testing.T) {
	// Test tilde alone
	result := expandPath("~")
	if result == "~" {
		t.Error("expandPath(~) should expand to home directory")
	}

	// Test tilde with path
	result = expandPath("~/projects")
	if result == "~/projects" {
		t.Error("expandPath(~/projects) should expand to home/projects")
	}
}

// TestShortenPath tests the path shortening helper
func TestShortenPath(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "empty", input: "", want: ""},
		{name: "tmp path", input: "/tmp/test", want: "/tmp/test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shortenPath(tt.input)
			if got != tt.want {
				t.Errorf("shortenPath(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestSanitizeSessionName tests the session name sanitization
func TestSanitizeSessionName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "simple", input: "myproject", want: "myproject"},
		{name: "uppercase", input: "MyProject", want: "myproject"},
		{name: "spaces", input: "my project", want: "my-project"},
		{name: "underscores", input: "my_project", want: "my-project"},
		{name: "multiple dashes", input: "my--project", want: "my-project"},
		{name: "special chars", input: "my@project#123", want: "myproject123"},
		{name: "leading dash", input: "-myproject", want: "myproject"},
		{name: "trailing dash", input: "myproject-", want: "myproject"},
		{name: "empty", input: "", want: "session"},
		{name: "only special", input: "@#$%", want: "session"},
		{name: "whitespace only", input: "   ", want: "session"},
		{name: "mixed", input: " My-Project_123 ", want: "my-project-123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeSessionName(tt.input)
			if got != tt.want {
				t.Errorf("sanitizeSessionName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestProjectListItem tests the Project type's list.Item implementation
func TestProjectListItem(t *testing.T) {
	p := Project{
		Name:    "Test Project",
		Session: "test-session",
		Path:    "/home/user/projects/test",
		Layout:  "dev-3",
		Status:  StatusRunning,
	}

	// Test Title
	title := p.Title()
	if title == "" {
		t.Error("Project.Title() should not be empty")
	}
	if title != "● Test Project" {
		t.Errorf("Project.Title() = %q, want %q", title, "● Test Project")
	}

	// Test Description
	desc := p.Description()
	if desc == "" {
		t.Error("Project.Description() should not be empty")
	}

	// Test FilterValue
	filter := p.FilterValue()
	if filter != "Test Project" {
		t.Errorf("Project.FilterValue() = %q, want %q", filter, "Test Project")
	}
}

// TestProjectDescriptionEmpty tests description when path is empty
func TestProjectDescriptionEmpty(t *testing.T) {
	p := Project{
		Name:   "Test",
		Path:   "",
		Status: StatusStopped,
	}

	desc := p.Description()
	if desc != "No path configured" {
		t.Errorf("Project.Description() = %q, want %q", desc, "No path configured")
	}
}

// TestGitProjectListItem tests the GitProject type's list.Item implementation
func TestGitProjectListItem(t *testing.T) {
	gp := GitProject{
		Name: "my-repo",
		Path: "/home/user/projects/my-repo",
	}

	// Test Title
	title := gp.Title()
	if title != "📁 my-repo" {
		t.Errorf("GitProject.Title() = %q, want %q", title, "📁 my-repo")
	}

	// Test FilterValue
	filter := gp.FilterValue()
	if filter != "my-repo" {
		t.Errorf("GitProject.FilterValue() = %q, want %q", filter, "my-repo")
	}
}

// TestKeyBindings tests key binding creation
func TestKeyBindings(t *testing.T) {
	// Test delegate key map
	dk := newDelegateKeyMap()
	if dk == nil {
		t.Fatal("newDelegateKeyMap() returned nil")
	}

	// Verify choose binding
	if !key.Matches(tea.KeyMsg{Type: tea.KeyEnter}, dk.choose) {
		t.Error("choose binding should match Enter key")
	}

	// Verify kill binding
	if !key.Matches(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'K'}}, dk.kill) {
		t.Error("kill binding should match K key")
	}

	// Test list key map
	lk := newListKeyMap()
	if lk == nil {
		t.Fatal("newListKeyMap() returned nil")
	}

	// Verify refresh binding
	if !key.Matches(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}, lk.refresh) {
		t.Error("refresh binding should match r key")
	}

	// Verify toggle help binding
	if !key.Matches(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}, lk.toggleHelp) {
		t.Error("toggleHelp binding should match ? key")
	}
}

// TestDelegateKeyMapHelp tests ShortHelp and FullHelp implementations
func TestDelegateKeyMapHelp(t *testing.T) {
	dk := newDelegateKeyMap()

	// Test ShortHelp
	shortHelp := dk.ShortHelp()
	if len(shortHelp) != 2 {
		t.Errorf("ShortHelp() returned %d bindings, want 2", len(shortHelp))
	}

	// Test FullHelp
	fullHelp := dk.FullHelp()
	if len(fullHelp) != 1 {
		t.Errorf("FullHelp() returned %d groups, want 1", len(fullHelp))
	}
	if len(fullHelp[0]) != 2 {
		t.Errorf("FullHelp()[0] has %d bindings, want 2", len(fullHelp[0]))
	}
}

// TestViewStateConstants tests view state enum values
func TestViewStateConstants(t *testing.T) {
	// Ensure distinct values
	states := map[ViewState]string{
		StateHome:          "home",
		StateProjectPicker: "picker",
		StateConfirmKill:   "confirm",
	}

	seen := make(map[ViewState]bool)
	for state := range states {
		if seen[state] {
			t.Errorf("ViewState %d is duplicated", state)
		}
		seen[state] = true
	}
}

// TestStatusConstants tests status enum values
func TestStatusConstants(t *testing.T) {
	// Ensure distinct values
	statuses := map[Status]string{
		StatusStopped: "stopped",
		StatusRunning: "running",
		StatusCurrent: "current",
	}

	seen := make(map[Status]bool)
	for status := range statuses {
		if seen[status] {
			t.Errorf("Status %d is duplicated", status)
		}
		seen[status] = true
	}
}
