# Review: Go TUI + CLI Best Practices (2025)

**Reviewed file:** `docs/BEST-PRACTICES.md`  
**Review date:** November 26, 2025

---

## Overall Assessment

⭐⭐⭐⭐ **Strong document** — Well-structured, opinionated, and practical. Covers the essential ground for building production-quality Go TUI applications with the Charm stack.

---

## Strengths

### 1. Clear Technology Stack
The document commits to the Charm ecosystem (`bubbletea`, `bubbles`, `lipgloss`, `huh`, `wish`) rather than being framework-agnostic. This is the right call — the Charm stack is the de facto standard for Go TUIs in 2025.

### 2. Correct Architectural Principles
- **Separation of concerns** — Business logic stays in plain Go packages, TUI is a presentation layer.
- **Pure Update functions** — No blocking IO in `Update`, use `tea.Cmd` for async work.
- **Small models** — State only, no IO handles or loggers.

These are the fundamentals that prevent spaghetti TUI code.

### 3. Practical UX Guidance
Section 5 (UX Rules) hits the right notes:
- Footer keybind hints
- `?` for help overlay
- Progress indicators with context (`X of Y`)
- `Enter` confirms, `Esc` cancels convention

### 4. Testing Strategy
The testing section correctly identifies:
- Unit test domain logic separately
- Test `Update` as a pure function
- Use `teatest` for integration tests
- Don't depend on terminal size in tests

---

## Gaps & Suggestions

### Missing: Error Handling Patterns
No guidance on how to surface errors in a TUI context:
- When to show inline errors vs modal dialogs
- How to structure error messages for `tea.Msg`
- Graceful degradation when background work fails

**Suggestion:** Add a dedicated section or expand Section 6 with error handling patterns.

### Missing: State Management at Scale
For larger apps with multiple screens/views:
- No mention of navigation patterns (stack-based, tab-based)
- No guidance on shared state vs screen-local state
- No discussion of when to split into multiple `tea.Program` instances

**Suggestion:** Add a subsection in Section 3 covering multi-screen navigation.

### Missing: Logging & Debugging in TUI Context
Section 6 briefly mentions `DEBUG` env var logging, but TUI debugging is notoriously tricky (stdout is the UI). Worth expanding:
- `bubbletea.WithOutput(io.Discard)` for test runs
- Logging to a file or `tea.Println` for debug output
- Using `tea.Sequence` for debugging command chains

### Missing: Configuration File Handling
Section 10 mentions `$XDG_CONFIG_HOME` but no guidance on:
- Config file format (YAML, TOML, JSON)
- How to load config and pass to TUI
- Live config reloading patterns

### Weak: SSH TUI Section
Section 9 covers security essentials but lacks practical implementation details:
- No mention of `wish.Middleware` patterns
- No guidance on session management
- No example of auth flow

This section feels more like a security checklist than a best practices guide.

---

## Nitpicks

| Line | Issue |
|------|-------|
| 7 | Consider mentioning `charmbracelet/log` for styled logging |
| 18 | `taskfile` → clarify this refers to `go-task/task` |
| 54 | List is incomplete — `spinner`, `textinput`, `table` are common Bubbles too |
| 94 | `teatest` link or import path would help |
| 121 | `entr` example is good, but `air` or `watchexec` are also popular |

---

## Applicability to PeakyPanes

This project (`peakypanes2`) appears to be a tmux session manager with a TUI. Reviewing against the best practices:

| Practice | Status in Codebase |
|----------|-------------------|
| `cmd/app/main.go` layout | ✅ Uses `cmd/peakypanes/main.go` |
| Separate TUI in `internal/tui/` | ✅ Present |
| Business logic separation | ✅ `internal/layout/`, `internal/tmuxctl/` |
| Theme centralization | ⚠️ Not visible — check if styles are scattered |
| `teatest` or similar | ❓ No test files visible for TUI |

**Recommendation:** Add TUI unit tests following Section 8 guidelines.

---

## Verdict

This is a **solid reference document** for the team. It correctly captures modern Go TUI patterns and avoids common pitfalls.

**Priority improvements:**
1. Add error handling patterns section
2. Expand state management for multi-screen apps
3. Add concrete examples/code snippets for key patterns

Would elevate this from "good checklist" to "comprehensive guide."
