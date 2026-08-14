<div align="center">
<h1>peky</h1>

**All your terminal AI agents, just one app.**

Peky puts terminal agents such as Claude Code, Codex CLI, pi, and opencode into a single app to keep projects organized.

[![homebrew](https://img.shields.io/badge/homebrew-regenrek%2Ftap%2Fpeky-2e933c?logo=homebrew)](https://github.com/regenrek/homebrew-tap)
[![pkg.go.dev](https://pkg.go.dev/badge/github.com/regenrek/peakypanes.svg)](https://pkg.go.dev/github.com/regenrek/peakypanes)

</div>

## Features

- 🧠 **AI agent orchestration** - Run agents side by side with quick replies, slash commands, and broadcast
- 🗂️ **Multi-project dashboard** - See projects and sessions in one TUI
- 🖱️ **Mouse support** - Create, select, resize, and drag panes
- 🎯 **Zero config** - Run `peky` in any directory
- 🧭 **Persistent native daemon** - Sessions keep running after the UI exits
- 📜 **Scrollback and copy mode** - Navigate output and copy from native panes
- 📁 **Project-local config** - Commit `.peky.yml` to git for team sharing

## Quick Start

### Install

The supported packaged installation is Homebrew:

```bash
brew tap regenrek/tap
brew install regenrek/tap/peky
brew services start peky
peky
```

Windows users should install a GitHub release or build from source.

### Usage

`peky` opens the dashboard. Start a supported two-pane session and open its
dashboard with:

```bash
peky start --layout split-v
```

`peky start` uses the built-in `auto` definition; see the split-layout
limitation below.

Create project-local configuration for a team:

```bash
cd your-project
peky init --local
# Edit .peky.yml
git add .peky.yml
```

## Built-in Layouts

The embedded layouts are:

| Layout | Description |
|--------|-------------|
| `auto` | Three-pane default definition; see the split-layout limitation below |
| `split-v` | Two vertical panes (left/right) |
| `split-h` | Two horizontal panes (top/bottom) |
| `3x3` | 3x3 grid (9 panes) |
| `4x3` | 4x3 grid (12 panes) |

`tauri-debug`, `fullstack`, and `go-dev` are not built-in layouts. Names such as
`tauri-debug` may be used for custom layouts, including the example in the
[Layout Builder Guide](docs/layout-builder.md).

The current split-tree builder supports one or two ordered panes. The embedded
`auto` definition contains three panes and therefore cannot currently start
successfully; use a grid layout such as `3x3` or `4x3` for more panes.

List or export layouts:

```bash
peky layouts
peky layouts export 3x3
```

`layouts export` prints a standalone layout. It can be saved directly as
`.peky.yml`; the loader also accepts the wrapped form shown below.

## Configuration

Project-local configuration lives in `.peky.yml`; global configuration lives in
`~/.config/peky/config.yml`.

```yaml
# .peky.yml
session: my-project

layout:
  panes:
    - title: editor
      cmd: "${EDITOR:-}"
    - title: server
      cmd: "npm run dev"
      split: horizontal
```

A project-local file may also contain a standalone layout at the top level when
it has `panes` or `grid`; this is the form emitted by `peky layouts export`.
The `session` value overrides the default session name, which is a sanitized
version of the project directory name. `peky start --session NAME` takes
precedence over it.

See the [configuration reference](docs/configuration.md) and [Layout Builder
Guide](docs/layout-builder.md) for the full schema.

## Variable Expansion

Use `${PROJECT_PATH}`, `${PROJECT_NAME}`, environment variables, and
`${VAR:-default}` in layout commands. Nested braced expressions are not
expanded recursively; after the braced pass, `$HOME` and a leading `~/` are
expanded as conveniences.

```yaml
layout:
  vars:
    log_file: "/tmp/peky.log"
  panes:
    - cmd: "tail -f ${PROJECT_PATH}/logs/${PROJECT_NAME}.log"
    - cmd: "tail -f ${log_file}"
    - cmd: "${EDITOR:-vim}"
```

Use built-in and environment variables directly in the command when they need
to expand. Braced variables are supported; `$HOME` and a leading `~/` are also
expanded.

## CLI Overview

The command is `peky`. `peky` with no command opens the dashboard; `start`,
`open`, and `o` start a session and open the dashboard. A layout name or `.`
can be used as a start shorthand.

| Command | Purpose |
|---------|---------|
| `dashboard`, `ui` | Open the dashboard |
| `start`, `open`, `o` | Start a session and open the dashboard |
| `daemon [start\|stop\|restart]` | Manage the native daemon |
| `init` | Create global or project-local configuration |
| `layouts [export NAME]` | List or export layouts |
| `workspace [list\|open\|close\|close-all]` | Manage projects |
| `clone`, `c` | Clone a repository and start a session |
| `session ...` | List and manage sessions |
| `pane ...` | List and manage panes, input, and layout operations |
| `relay ...` | Create and manage pane relays |
| `events ...` | Watch or replay events |
| `context pack` | Create a context bundle |
| `nl [plan\|run]` | Plan or run a natural-language command |
| `version` | Show the version |
| `help` | Show command help |

Global flags are available on every command:

```text
--json             Emit JSON output
--timeout DURATION Override the command timeout
--yes, -y          Skip confirmations for side-effect commands
--version, -v      Show the version and exit
--fresh-config     Ignore global config and saved state
--temporary-run    Use temporary runtime and config directories
```

For command-specific arguments and flags, see the [CLI reference](docs/cli.md).

## Layout Detection

For `peky start` and its aliases, layout selection is evaluated in this order:

1. An explicit `--layout NAME` flag or layout-name shorthand. A named global
   layout overrides a built-in layout with the same name.
2. The project-local `.peky.yml` or `.peky.yaml` layout.
3. The built-in `auto` layout.

Global `projects` entries are used by the dashboard/workspace UI; they are not
matched automatically by a direct `peky start`. A global layout named `auto` is
selected only when requested explicitly; it does not replace the fallback.

## For Teams

1. Run `peky init --local` in your project.
2. Customize `.peky.yml` for your stack.
3. Commit it to git.
4. Teammates install peky and run `peky start --layout split-v` (or a grid layout).

## License

MIT
