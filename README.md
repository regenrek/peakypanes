# 🎩 Peaky Panes

```
████    █████    ███    █   █   █   █    ████      ███    █   █    █████    ████
█████   ████    █████   ████     ███     █████    █████   ███ █    ████    ████ 
█       █████   █   █   █  ██     █      █        █   █   █  ██    █████   █████
```

**Tmux layout manager with YAML-based configuration.**

Define your tmux layouts in YAML, share them with your team via git, and get consistent development environments everywhere.

## Features

- 📦 **Built-in layouts** - Works out of the box with sensible defaults
- 📁 **Project-local config** - Commit `.peakypanes.yml` to git for team sharing
- 🏠 **Global config** - Define layouts once, use everywhere
- 🔄 **Variable expansion** - Use `${EDITOR}`, `${PROJECT_PATH}`, etc.
- 🖥️ **Project-manager TUI** - Run `peakypanes` with no command to manage projects
- 🎯 **Zero config sessions** - Run `peakypanes start` with the built-in default
- ⚙️ **Session-scoped tmux options** - Configure tmux per-session without affecting global config

## Quick Start

### Prerequisites

- [tmux](https://github.com/tmux/tmux) must be installed and available on `PATH`.
- Go 1.23.5 or newer is required when building from source.

### Install

The canonical installation path is currently unresolved in this checkout: `go.mod`
declares `github.com/kregenrek/tmuxman`, while the repository remote is
`github.com/regenrek/peakypanes`. Confirm the published module path before using
`go install`.

### Usage

**Open the project manager:**
```bash
cd your-project
peakypanes
```

With no command, `peakypanes` opens the interactive project-manager TUI. It does
not immediately start the current directory's session. The TUI also requires
tmux.

**Start or attach to the current directory:**
```bash
peakypanes start
peakypanes open
```

**Use a specific layout:**
```bash
peakypanes start --layout dev-3
peakypanes start --layout fullstack
```

**Create project-local config (recommended for teams):**
```bash
cd your-project
peakypanes init --local
# Edit .peakypanes.yml
git add .peakypanes.yml  # Share with team
```

## Built-in Layouts

| Layout | Description |
|--------|-------------|
| `simple` | Single window with one pane |
| `dev-2` | Simple two-pane layout: editor left, terminal right |
| `dev-3` | Three-pane layout: editor and two additional panes (default) |
| `fullstack` | Fullstack development: editor, dev server, and logs |
| `go-dev` | Go development: editor, run, and tests |
| `split-h` | Two horizontal panes (top/bottom) |
| `split-v` | Two vertical panes (left/right) |

`tauri-debug` is not a built-in layout. The layout-builder guide includes it as
an example of a custom layout.

```bash
# List all layouts
peakypanes layouts

# Export a layout to customize
peakypanes layouts export dev-3
```

## Configuration

> 📖 **[Layout Builder Guide](docs/layout-builder.md)** - Detailed documentation on creating custom layouts, pane arrangements, and tmux options.

### Project-Local (`.peakypanes.yml`)

Create in your project root for team-shared layouts:

```yaml
# .peakypanes.yml
layout:
  windows:
    - name: dev
      panes:
        - title: editor
          cmd: "${EDITOR:-}"
        - title: server
          cmd: "npm run dev"
          split: horizontal
          size: "40%"
        - title: shell
          cmd: ""
          split: vertical

    - name: logs
      panes:
        - title: docker
          cmd: "docker compose logs -f"
```

The session name defaults to the project directory name; use `start --session`
to override it.

### Global Config (`~/.config/peakypanes/config.yml`)

For personal layouts and multi-project management:

```yaml
# Global settings
tmux:
  config: ~/.config/tmux/tmux.conf

# Custom layouts
layouts:
  my-custom:
    windows:
      - name: main
        panes:
          - title: code
            cmd: nvim
          - title: term
            cmd: ""

# Projects shown in the project-manager TUI
projects:
  - name: webapp
    session: webapp
    path: ~/projects/webapp
    layout: fullstack
```

The `projects` list is used by the no-argument project-manager TUI. It is not
part of automatic layout detection for `start` or `open`.

## Variable Expansion

Use variables in your layouts:

| Variable | Description |
|----------|-------------|
| `${PROJECT_PATH}` | Project path used by `start` |
| `${PROJECT_NAME}` | Directory name |
| `${EDITOR}` | Your $EDITOR |
| `${VAR:-default}` | Env var with default |

Expansion is single-pass: variables inside a `vars` value are not expanded
again after that value is substituted. Put built-in and environment variables
in the command that uses them when they need expansion:

```yaml
layout:
  windows:
    - name: dev
      panes:
        - cmd: "tail -f ${PROJECT_PATH}/logs/${PROJECT_NAME}.log"
        - cmd: "${EDITOR:-vim}"
```

## Commands

| Command | Aliases and arguments | Purpose |
|---------|------------------------|---------|
| `peakypanes` | — | Open the interactive project-manager TUI |
| `peakypanes start` | `-l, --layout <name>`; `-s, --session <name>`; `-p, --path <dir>` | Start or attach a session |
| `peakypanes open` | `peakypanes o` | Alias for `start` |
| `peakypanes kill` | `peakypanes k [session-name]` | Kill a tmux session |
| `peakypanes init` | `-l, --local`; `--layout <name>`; `-f, --force` | Create global or project-local configuration |
| `peakypanes layouts` | `layouts export <name>` | List layouts or export YAML |
| `peakypanes clone` | `peakypanes c <url\|user/repo>` | Clone to `~/projects/<repo>` and start a session |
| `peakypanes version` | `-v, --version` | Show the version |
| `peakypanes help` | `-h, --help` | Show top-level help |

The `start`, `kill`, `init`, and `layouts` commands expose `-h` and `--help`.
A non-flag argument in `start` is also accepted as a layout-name shortcut, for
example `peakypanes dev-3`. Run `peakypanes <command> --help` for
command-specific options.

## How Layout Detection Works

For `start` and `open`, layout selection is evaluated in this order:

1. The `--layout` flag or a positional layout-name shortcut.
2. `.peakypanes.yml` or `.peakypanes.yaml` in the project directory.
3. The `dev-3` default.

Global layouts from `~/.config/peakypanes/layouts/` and the `layouts` section of
`~/.config/peakypanes/config.yml` are available by name. A global layout with a
given name overrides the embedded layout with that name, including `dev-3`,
when the name is selected explicitly or as the default.

## For Teams

1. Run `peakypanes init --local` in your project
2. Customize `.peakypanes.yml` for your stack
3. Commit to git
4. Teammates install peakypanes and run `peakypanes start` - done!

## License

MIT
