# Layout Builder Guide

This guide covers the layout YAML used by peky. Layouts can use ordered pane
splits or exact rectangular grids, plus optional command and input automation.

## Table of Contents

- [Basic Structure](#basic-structure)
- [Pane Layouts](#pane-layouts)
- [Split Directions](#split-directions)
- [Variables](#variables)
- [Automation Sends](#automation-sends)
- [Examples](#examples)
- [Configuration Precedence](#configuration-precedence)
- [Tips](#tips)

---

## Basic Structure

A project-local `.peky.yml` normally wraps the layout in `layout:`:

```yaml
# Optional: defaults to a sanitized project directory name
session: my-project

layout:
  name: my-layout
  description: "Description of what this layout is for"

  vars:
    log_file: "/tmp/peky.log"

  panes:
    - title: editor
      cmd: "${EDITOR:-}"
    - title: server
      cmd: "npm run dev"
      split: horizontal
```

The standalone layout printed by `peky layouts export NAME` may also be saved
as `.peky.yml`. The loader treats a top-level file containing `panes` or `grid`
as the project layout. `peky start --session NAME` overrides the `session` value.

`settings.width` and `settings.height` are accepted by the YAML schema but are
not currently applied when a session starts. Do not use them to resize panes.
Likewise, pane `setup` is parsed but not run, and pane `enabled` is parsed but
not evaluated. Do not rely on those fields for runtime behavior.

---

## Pane Layouts

### Split Layouts

Use `panes` for an ordered layout. The first pane is the base pane and the
second pane can split it. The current split-tree builder does not support a
third ordered pane; use `grid` for layouts with more than two panes.

```yaml
layout:
  panes:
    - title: editor
      cmd: "${EDITOR:-vim}"
    - title: server
      cmd: "npm run dev"
      split: horizontal
```

### Exact Grid Layouts

Use `grid` when you need predictable rows and columns:

```yaml
layout:
  grid: 2x2
  commands:
    - "claude"
    - "npm run dev"
    - "tail -f app.log"
    - ""
  titles:
    - agent
    - dev
    - logs
    - shell
```

The `commands` and `titles` arrays are row-major. A top-level `command` is used
for every grid pane. If `commands` is shorter than the grid, its remaining
panes use `command` when it is set, otherwise they start empty. `panes` can
override the command and title for individual grid positions.

Built-in grid layouts are `3x3` and `4x3`. The `--panes N` flag on
`peky session start` can generate a balanced grid without a layout name; it
cannot be combined with `--layout`.

---

## Split Directions

For split layouts, use `split` on panes after the first:

| Split | Result |
|-------|--------|
| `horizontal` | Creates left/right panes |
| `vertical` | Creates top/bottom panes |

`v` is also accepted for vertical splits. Other values use the horizontal
split behavior.

### Size Control

Use `size` on a pane after the first to set the percentage of the split given to
the new pane:

```yaml
layout:
  panes:
    - title: main
      cmd: ""
    - title: side
      cmd: ""
      split: horizontal
      size: "30%"
```

A `size` value on the first pane is ignored. Grid layouts calculate equal row
and column sizes instead of using pane split sizes.

---

## Variables

### Special and Environment Variables

| Variable | Description |
|----------|-------------|
| `${PROJECT_PATH}` | Absolute project path used by `start` |
| `${PROJECT_NAME}` | Project directory name |
| `${HOME}` | The `HOME` environment variable |
| `${VAR:-default}` | Environment variable with a fallback |

Use braced variables in layout values:

```yaml
layout:
  panes:
    - title: editor
      cmd: "${EDITOR:-vim}"
    - title: logs
      cmd: "tail -f ${PROJECT_PATH}/logs/${PROJECT_NAME}.log"
```

Nested braced expressions are not expanded recursively. After the braced
pass, `$HOME` and a leading `~/` are expanded as conveniences. Use direct
expansion for special or environment variables that must be resolved:

```yaml
layout:
  vars:
    log_file: "/tmp/peky.log"
  panes:
    - cmd: "tail -f ${log_file}"
    - cmd: "${HOME}/bin/agent"
```

Custom variables are checked first, then the environment, then the `:-default`
value. `$HOME` and a leading `~/` are expanded as additional conveniences.
Bare `$EDITOR` is not part of the braced variable syntax; use `${EDITOR}` or a
default expression.

---

## Automation Sends

Use `broadcast_send` to send an input action to every pane after startup, or
`direct_send` on a pane to target that pane. Each action accepts `text`,
`send_delay_ms`, `wait_for_output`, `submit`, and `submit_delay_ms`.

```yaml
layout:
  command: "claude"
  broadcast_send:
    - text: "give me a status update"
      wait_for_output: true
      submit: true

  panes:
    - title: first
      cmd: "claude"
      direct_send:
        - text: "work on the first task"
          send_delay_ms: 1000
          submit: true
```

A trailing newline is added automatically unless `submit: true` is used;
`submit: true` sends Enter separately. If `send_delay_ms` is omitted, the
first output is awaited up to the default delay.

`session_restore: true`, `false`, or `private` is accepted on a pane, but the
current normal startup expansion does not preserve this per-pane override. Do
not rely on it until that path is fixed. The global `session_restore`
configuration controls daemon snapshot behavior separately.

---

## Examples

### Tauri/Rust Development (custom example)

`tauri-debug` is an example layout name, not a built-in layout. For a
project-local `.peky.yml`, keep the `layout:` wrapper. For a standalone layout
file, remove that wrapper if you want the file to contain only the layout.

```yaml
layout:
  name: tauri-debug
  description: "Tauri development"
  grid: 2x2
  commands:
    - "codex"
    - "bun dev:tauri"
    - "tail -F ${HOME}/Library/Logs/${PROJECT_NAME}/rust.log"
    - ""
  titles:
    - codex
    - bun-dev
    - rust-logs
    - shell
```

### Simple Two-Pane Layout

```yaml
layout:
  panes:
    - title: editor
      cmd: "${EDITOR:-vim}"
    - title: terminal
      cmd: ""
      split: horizontal
      size: "40%"
```

---

## Configuration Precedence

For `peky start`, `peky open`, and `peky o`, layout selection is evaluated in
this order:

1. An explicit `--layout NAME` flag or layout-name shorthand.
2. The project-local `.peky.yml` or `.peky.yaml` layout.
3. The built-in `auto` layout.

For an explicit name, a global layout from the default layouts directory or
the global config's `layouts` map overrides a built-in layout with the same
name. A global layout named `auto` is selected only when requested explicitly;
it does not replace the fallback. Entries in the global `projects` list are
used by the dashboard/workspace UI and are not matched automatically by a
direct `start`.

---

## Tips

### Prefer `grid:` for Exact Grids

Use `grid: 2x3`, `grid: 3x3`, or another valid rectangular grid when pane
positions should be predictable. Use ordered `panes` when split direction and
per-pane sizes matter.

### Empty Commands

Use `cmd: ""` for a pane that should open as an interactive shell.

### Review Commands

Layout commands are executed when the session starts. Review a project-local
`.peky.yml` before starting an untrusted repository.
