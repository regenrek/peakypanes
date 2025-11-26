# Go TUI + CLI Rules (2025)

## 0. Stack

- MUST use a clear stack:
  - CLI args: `spf13/cobra`, `urfave/cli/v2`, or `ffcli` for multi-command tools.
  - TUI: Charm stack: `bubbletea` (core), `bubbles` (widgets), `lipgloss` (layout/styling), `huh` (forms).
  - SSH TUIs: `wish` for apps over SSH when needed.
  - Logging: `charmbracelet/log` for styled terminal logging that respects TUI context.
- SHOULD keep business logic in plain Go packages with no TUI imports. TUI is just another front end.
- MUST be able to run core logic without a terminal (unit tests, headless runs).

## 1. Project Layout

- MUST follow a standard Go layout for CLIs:
  - `cmd/app/main.go` for entrypoint.
  - `internal/` or top level packages for domain logic.
  - `internal/tui/` for Bubble Tea models, `internal/ui/theme` for Lip Gloss styles.
- SHOULD have a `Makefile` or `taskfile` ([go-task/task](https://taskfile.dev)) with `test`, `lint`, `build`, `dev` targets.
- SHOULD keep each TUI model in its own file; large screens in their own directory.

## 2. CLI Interface Rules (non-TUI)

- MUST have a non-interactive mode for automation:
  - Respect `--help`, `--version`, `--json`, `--silent`.
- MUST detect TTY vs non-TTY:
  - If not a TTY, default to plain CLI output, not a full-screen TUI.
- SHOULD:
  - Keep subcommand tree shallow and predictable (`app list`, `app show`, `app edit`, `app delete`).
  - Provide shell completion if using Cobra or urfave/cli.

## 3. TUI Architecture (Bubble Tea)

- MUST keep `tea.Model` small:
  - State only, no IO handles, no loggers.
  - One model per screen or widget; compose into a tree of models.
- MUST keep `Update` fast:
  - No blocking IO inside `Update`.
  - Use `tea.Cmd` for IO and long work, send results back as messages.
- SHOULD:
  - Treat `Msg` as domain events (`LoadedConfigMsg`, `JobProgressMsg`, not `string`).
  - Use custom messages for async results, errors, ticks.
- MUST handle `tea.WindowSizeMsg` correctly and store width/height in the model for responsive layout.

### Multi-Screen Navigation

- MUST use a clear navigation pattern for apps with multiple screens:
  - **Stack-based:** Push/pop screens like a navigation stack. Good for drill-down flows.
  - **Tab-based:** Switch between peer screens. Good for dashboards.
  - **State machine:** Explicit states and transitions. Good for wizard flows.
- SHOULD:
  - Keep a `currentScreen` enum or interface in your root model.
  - Pass shared state down to child models; bubble events up via messages.
  - Avoid global state; prefer explicit dependency injection through model fields.

```go
type Screen int
const (
    ScreenList Screen = iota
    ScreenDetail
    ScreenEdit
)

type Model struct {
    screen      Screen
    listModel   list.Model
    detailModel DetailModel
    // shared state
    selectedID  string
}
```

## 4. Layout & Styling (Lip Gloss + Bubbles)

- MUST:
  - Centralize theme in one package (colors, borders, padding).
  - Use `lipgloss.JoinHorizontal` and `JoinVertical` instead of hand-written spacing.
  - Handle small terminals by collapsing sections or truncating, not by breaking rendering.
- SHOULD:
  - Support `NO_COLOR` and a `--no-color` flag, use `lipgloss.NoColor` where needed.
  - Keep layout math simple and avoid deep nesting that makes width and height reasoning fragile.
- MUST use Bubbles components where they fit:
  - `list`, `table`, `textarea`, `textinput`, `spinner`, `progress`, `viewport`, etc, instead of rolling your own widget every time.

## 5. UX Rules For TUIs

- MUST provide inline help:
  - Footer bar with keys: `↑↓ / j k: move`, `enter: select`, `q: quit`, `?: help`.
  - `?` opens a help view or overlay with key cheatsheet.
- MUST use clear progress feedback:
  - Spinner for quick unknown work.
  - `X of Y` for multi-step operations.
  - Progress bar for long tasks or downloads.
- SHOULD:
  - Always show current context in a header (active resource, path, filter).
  - Use obvious defaults: `Enter` confirms, `Esc` cancels.
  - Keep forms linear: one focusable input at a time unless it is a dashboard.

## 6. Async Work & Performance

- MUST:
  - Run heavy tasks in goroutines, return a `tea.Cmd` that waits and sends a result message.
  - Use `context.Context` for cancellation so quitting the TUI cancels background work.
- SHOULD:
  - Batch frequent updates (throttle progress to a few times per second).
  - Log debug info to a file (message traces) when a `DEBUG` env var is set.

### Logging & Debugging

- MUST:
  - Never write to stdout/stderr directly in TUI mode — it corrupts the screen.
  - Use file-based logging during development: `log.SetOutput(file)` or `charmbracelet/log` with file output.
- SHOULD:
  - Use `tea.WithOutput(io.Discard)` in tests to suppress TUI rendering.
  - Enable verbose logging with `DEBUG=1` or `--debug` flag.
  - Log message types in `Update` to trace state transitions.

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if os.Getenv("DEBUG") != "" {
        log.Printf("msg: %T %+v", msg, msg)
    }
    // ...
}
```

## 7. Error Handling

- MUST:
  - Define typed error messages: `type ErrMsg struct { Err error }`.
  - Handle errors in `Update` and update model state to reflect the error.
  - Always provide a way for the user to recover or dismiss errors.
- SHOULD:
  - Show transient errors inline (toast/status bar) for non-blocking issues.
  - Show modal errors for blocking issues that require user action.
  - Include context in error messages: what failed and what the user can do.

```go
type ErrMsg struct {
    Err     error
    Context string // e.g., "loading config"
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case ErrMsg:
        m.err = msg.Err
        m.errContext = msg.Context
        return m, nil
    }
    // ...
}
```

- MUST:
  - Clear error state when the user takes corrective action or dismisses.
  - Not swallow errors silently — always surface them to the user or log them.

## 8. Forms, Prompts, And Input

- MUST:
  - Use `huh` or Bubbles input components for multi-step prompts rather than ad hoc input loops.
  - Validate user input early and show clear inline errors.
- SHOULD:
  - Provide sensible defaults and remember previous values in config.
  - For destructive actions, use a confirmation step with a distinct style (red border, clear wording).

## 9. Testing TUIs

- MUST:
  - Unit test the core domain logic without Bubble Tea.
  - Unit test `Update` as a pure function: given `model` plus `Msg` you assert new state and outgoing command.
- SHOULD:
  - Use [`teatest`](https://github.com/charmbracelet/x/tree/main/exp/teatest) or similar to run the full TUI in a virtual terminal and assert view output and state.
  - For complex widgets, snapshot-test the rendered `View()` output (with a golden file or a helper library).
- MUST ensure tests do not depend on terminal size or color; use consistent settings in tests.

## 10. SSH TUIs (Wish)

- MUST:
  - Treat SSH input as untrusted. No shell escapes, no direct file access without checks.
  - Lock down SSH auth methods and allowed keys.
  - Separate the SSH server process from your core library so it is easy to harden and monitor.
- SHOULD:
  - Detect client terminal features and degrade nicely (no bold colors on very basic terminals).
  - Log failed auth, connection counts, and panics.

### Implementation Patterns

- MUST use `wish.Middleware` to compose functionality:
  - Auth middleware → logging middleware → your TUI handler.
- SHOULD:
  - Store per-session state in the `ssh.Session` context.
  - Use `wish.WithHostKeyPath` for persistent host keys.
  - Implement session timeouts for idle connections.

```go
func main() {
    s, err := wish.NewServer(
        wish.WithAddress(":2222"),
        wish.WithHostKeyPath(".ssh/term_info_ed25519"),
        wish.WithMiddleware(
            bubbletea.Middleware(teaHandler),
            activeterm.Middleware(),
            logging.Middleware(),
        ),
    )
    // ...
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
    pty, _, _ := s.Pty()
    return NewModel(pty.Window.Width, pty.Window.Height), nil
}
```

## 11. Configuration

- MUST:
  - Support a config file location following platform conventions:
    - Linux/macOS: `$XDG_CONFIG_HOME/app/config.yaml` or `~/.config/app/config.yaml`
    - Windows: `%APPDATA%\app\config.yaml`
  - Allow environment variable overrides for all config values.
  - Document all config options with examples.
- SHOULD:
  - Use YAML or TOML for human-editable config (avoid JSON for config files).
  - Use [`viper`](https://github.com/spf13/viper) or [`koanf`](https://github.com/knadh/koanf) for config loading.
  - Provide a `--config` flag to specify an alternate config file.
  - Create default config on first run with helpful comments.

```go
func LoadConfig() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("$XDG_CONFIG_HOME/myapp")
    viper.AddConfigPath("$HOME/.config/myapp")
    viper.AutomaticEnv()
    viper.SetEnvPrefix("MYAPP")
    
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return nil, err
        }
    }
    // ...
}
```

## 12. Distribution & Tooling

- MUST:
  - Cross-compile with `goreleaser` or similar for macOS, Linux, Windows.
  - Ship checksums and a changelog per release.
- SHOULD:
  - Provide prebuilt binaries plus `go install` instructions.
  - Support environment overrides for CI and scripts.

## 13. Developer Experience

- MUST:
  - Provide a fast dev loop: `make dev` that runs the TUI and reloads on changes (using [`entr`](https://eradman.com/entrproject/), [`air`](https://github.com/air-verse/air), or [`watchexec`](https://github.com/watchexec/watchexec)).
  - Document key mappings and architecture in a short `README-DEV-TUI.md`.
- SHOULD:
  - Record short demo gifs or terminal recordings of key flows and keep them updated.
