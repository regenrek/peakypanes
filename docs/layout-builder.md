# Layout Builder Guide

This guide covers everything you need to know about creating custom layouts in Peaky Panes, including pane arrangements, tmux options, and advanced configuration.

## Table of Contents

- [Basic Structure](#basic-structure)
- [Pane Layouts](#pane-layouts)
- [Split Directions](#split-directions)
- [Tmux Options](#tmux-options)
- [Variables](#variables)
- [Multi-Window Layouts](#multi-window-layouts)
- [Examples](#examples)

---

## Basic Structure

A `.peakypanes.yml` file has this structure:

```yaml
# Session names default to the project directory; use start --session to override.
layout:
  name: my-layout
  description: "Description of what this layout is for"
  
  vars:
    # Custom literal variables
    log_file: "/tmp/peakypanes.log"
  
  settings:
    tmux_options:       # Session-scoped tmux options
      history-limit: "50000"
  
  windows:
    - name: dev
      layout: tiled     # tmux layout algorithm
      panes:
        - title: editor
          cmd: "${EDITOR:-}"
```

`width`, `height`, `bind_keys`, `setup`, and `enabled` are accepted in the YAML
schema but are not currently applied or evaluated when Peaky Panes creates a
session. Do not rely on them to resize terminals, create bindings, run setup
commands, or conditionally disable panes.

Commands in `cmd` fields are executed by tmux. Treat layout files as executable
configuration and review them before starting an untrusted project.

---

## Pane Layouts

### Automatic Layouts (Recommended)

Use the `layout` field on a window to let tmux arrange panes automatically:

| Layout | Description |
|--------|-------------|
| `tiled` | Equal-sized grid (best for 4 panes = 2x2) |
| `even-horizontal` | Side-by-side columns |
| `even-vertical` | Stacked rows |
| `main-horizontal` | Large pane on top, others below |
| `main-vertical` | Large pane on left, others on right |

```yaml
windows:
  - name: dev
    layout: tiled      # Automatic 2x2 grid with 4 panes
    panes:
      - title: top-left
        cmd: ""
      - title: top-right
        cmd: ""
      - title: bottom-left
        cmd: ""
      - title: bottom-right
        cmd: ""
```

### 2x2 Grid Example

The simplest way to get a 2x2 grid:

```yaml
windows:
  - name: main
    layout: tiled
    panes:
      - title: codex
        cmd: "codex"
      - title: dev
        cmd: "npm run dev"
      - title: logs
        cmd: "tail -f app.log"
      - title: shell
        cmd: ""
```

---

## Split Directions

For precise control, use `split` on individual panes:

| Split | Result |
|-------|--------|
| `horizontal` | Creates left/right panes |
| `vertical` | Creates top/bottom panes |

```yaml
windows:
  - name: dev
    panes:
      - title: editor        # First pane (no split)
        cmd: "${EDITOR:-}"
      - title: server        # Splits horizontally from editor
        cmd: "npm run dev"
        split: horizontal
        size: "40%"
      - title: shell         # Splits vertically from server
        cmd: ""
        split: vertical
```

This creates:
```
+------------------+----------+
|                  |  server  |
|     editor       +----------+
|                  |  shell   |
+------------------+----------+
```

### Size Control

Use `size` on a pane that is being split to pass its percentage to tmux. The
first pane is created before any split, so a `size` value on the first pane is
ignored.

```yaml
panes:
  - title: main
    cmd: ""
  - title: side
    cmd: ""
    split: horizontal
    size: "30%"           # New side pane gets 30% of the available space
```

---

## Tmux Options

Peaky Panes applies session-scoped tmux options that don't affect your global config.

### Default Options

These are applied automatically when Peaky Panes creates a new session. An
existing session is attached without applying the layout's options again:

| Option | Default | Description |
|--------|---------|-------------|
| `remain-on-exit` | `off` | Lets panes close normally after a command exits |

### Custom Options

Add your own tmux options per-layout:

```yaml
layout:
  settings:
    tmux_options:
      remain-on-exit: "on"       # Keep crashed panes visible
      history-limit: "50000"     # Scrollback buffer size
      mouse: "on"                # Enable mouse support
      status-position: "top"     # Status bar position
```

### Disabling Defaults

Override defaults if needed:

```yaml
layout:
  settings:
    tmux_options:
      remain-on-exit: "off"      # Let panes close on exit
```

### Key Bindings

`bind_keys` is part of the configuration schema, but it is not currently
applied to the tmux session. Configure bindings in tmux separately until
runtime support is added.

---

## Variables

### Special Variables

| Variable | Description |
|----------|-------------|
| `${PROJECT_PATH}` | Project path used by `start` |
| `${PROJECT_NAME}` | Project directory name |

### Environment Variables

Use any environment variable, such as `${HOME}`, with an optional default:

```yaml
panes:
  - title: editor
    cmd: "${EDITOR:-vim}"      # Use $EDITOR, fall back to vim
  - title: shell
    cmd: "${SHELL:-/bin/bash}" # Use $SHELL, fall back to bash
```

### Custom Variables

Define reusable literal values in the `vars` section:

```yaml
layout:
  vars:
    rust_log: "/tmp/peakypanes-rust.log"
    codex_log: "/tmp/peakypanes-codex.log"

  windows:
    - name: logs
      panes:
        - title: rust
          cmd: "tail -F ${rust_log}"
        - title: codex
          cmd: "tail -F ${codex_log}"
```

Expansion is single-pass. A variable's value is substituted as-is, so a value
such as `${HOME}/logs/${PROJECT_NAME}.log` will not expand the nested
expressions afterward. Put special or environment variables directly in the
command when they need expansion.

---

## Multi-Window Layouts

Create multiple tmux windows (tabs):

```yaml
layout:
  windows:
    # Window 1: Development
    - name: dev
      layout: tiled
      panes:
        - title: editor
          cmd: "${EDITOR:-}"
        - title: server
          cmd: "npm run dev"
        - title: test
          cmd: ""
        - title: shell
          cmd: ""
    
    # Window 2: Logs
    - name: logs
      layout: even-horizontal
      panes:
        - title: app-log
          cmd: "tail -f logs/app.log"
        - title: error-log
          cmd: "tail -f logs/error.log"
    
    # Window 3: Database
    - name: db
      panes:
        - title: psql
          cmd: "psql -d mydb"
```

---

## Examples

### Full-Stack Web Development

```yaml
layout:
  name: fullstack
  description: "Full-stack development with logs"
  
  settings:
    tmux_options:
      history-limit: "50000"
  
  windows:
    - name: code
      layout: main-vertical
      panes:
        - title: editor
          cmd: "${EDITOR:-}"
        - title: server
          cmd: "npm run dev"
        - title: shell
          cmd: ""
    
    - name: logs
      layout: even-horizontal
      panes:
        - title: frontend
          cmd: "tail -f logs/frontend.log"
        - title: backend
          cmd: "tail -f logs/backend.log"
```

### Tauri/Rust Development (custom example)

`tauri-debug` is an example name, not a built-in layout. For a project-local
`.peakypanes.yml`, keep the `layout:` wrapper below. For a standalone `.yml`
file in the global layouts directory, remove that wrapper. The `init --layout`
template option accepts built-in layout names only.

```yaml
layout:
  name: tauri-debug
  description: "Tauri development with codex agents"

  windows:
    - name: dev
      layout: tiled
      panes:
        - title: codex
          cmd: "RUST_LOG=debug codex"
        - title: bun-dev
          cmd: "bun dev:tauri"
        - title: codex-logs
          cmd: "tail -F ${HOME}/.spezi/codex/log/app-server.log | grep -Ev '\\bINFO\\b'"
        - title: rust-logs
          cmd: "tail -F ${HOME}/Library/Logs/${PROJECT_NAME}/rust.log"
```

### Go Development

```yaml
layout:
  name: go-dev
  
  windows:
    - name: dev
      panes:
        - title: editor
          cmd: "${EDITOR:-}"
        - title: run
          cmd: ""
          split: horizontal
          size: "40%"
        - title: test
          cmd: ""
          split: vertical
    
    - name: git
      panes:
        - title: lazygit
          cmd: "lazygit"
```

### Simple 2-Pane Layout

```yaml
layout:
  windows:
    - name: dev
      panes:
        - title: editor
          cmd: "${EDITOR:-}"
        - title: terminal
          cmd: ""
          split: horizontal
          size: "40%"
```

---

## Configuration Precedence

For `start` and `open`, layout selection is evaluated in this order:

1. `--layout` on the command line, or a positional layout-name shortcut.
2. `.peakypanes.yml` or `.peakypanes.yaml` in the project directory.
3. The `dev-3` default.

Global layouts in `~/.config/peakypanes/layouts/` and the `layouts` section of
`~/.config/peakypanes/config.yml` can be selected by name and override embedded
layouts with the same name. This includes a global `dev-3` overriding the
embedded default. Entries in the global `projects` list are for the
project-manager TUI; they are not consulted by `start` or `open` for automatic
layout selection.

---

## Tips

### Keep Crashed Panes Visible

Peaky Panes sets `remain-on-exit: off` by default, so panes close normally when commands exit. Set `remain-on-exit: on` in your layout's `tmux_options` if you want crashed commands to stay visible for debugging.

### Use `layout: tiled` for Grids

Don't manually calculate splits for grids. Use `layout: tiled` with the right number of panes:
- 2 panes → side by side
- 4 panes → 2x2 grid
- 6 panes → 2x3 or 3x2 grid

### Empty Commands

Use `cmd: ""` for panes where you want a shell ready for manual commands.

### Log Tailing with Fallbacks

Handle missing log files gracefully:

```yaml
cmd: "tail -F ${log_file} 2>/dev/null || echo 'Waiting for ${log_file}...'"
```

### Filter Noisy Logs

Reduce log noise with grep:

```yaml
cmd: "tail -F ${log_file} | grep -Ev 'DEBUG|TRACE|healthcheck'"
```
