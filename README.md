[简体中文](./README_ZH.md)

# pkmg

`pkmg` is a lightweight, local-first skill and script manager for humans and coding agents.

It helps you turn existing shell scripts into reusable local skills, inspect them with structured metadata, manage them in a minimal Web UI, and expose an agent-friendly command surface without turning the project into a heavy platform.

## Why pkmg

Most local automation starts as scattered scripts.

That works until:

- multiple skills need the same script
- scripts drift across projects
- agents need a safe way to discover and inspect what already exists
- editing and versioning local skills becomes messy

`pkmg` focuses on that local layer:

- manage reusable local skills under one place
- help agents list, search, inspect, and run existing skills
- provide a lightweight Web UI for creation, editing, copy, and version switching
- stay small: Go backend, static HTML, CDN-loaded frontend libraries

## Current Features

- Initialize a local skill workspace under `data/`
- Manage scripts in `data/scripts/`
- List and search skills with JSON output
- Inspect script metadata and content previews
- Run managed scripts
- Launch a lightweight local Web UI with:
  - local skill listing
  - search
  - create
  - edit
  - copy
  - version history
  - open containing folder

## Installation

### Go install

```bash
go install github.com/CoderSerio/pokemand-go@latest
```

### Local development build

```bash
git clone https://github.com/CoderSerio/pokemand-go.git
cd pokemand-go
go build -o bin/pkmg .
```

### Local global symlink for testing

On macOS, if `/opt/homebrew/bin` is already in your `PATH`:

```bash
go build -o bin/pkmg .
ln -sfn "$(pwd)/bin/pkmg" /opt/homebrew/bin/pkmg
pkmg --version
```

Remove the symlink later if needed:

```bash
rm /opt/homebrew/bin/pkmg
```

## Quick Start

Initialize the workspace:

```bash
pkmg init
```

Open or create a skill script:

```bash
pkmg open cleanup.sh
```

List skills:

```bash
pkmg list
pkmg list --json
```

Search skills:

```bash
pkmg search cleanup
pkmg search cleanup --json
```

Inspect a skill:

```bash
pkmg inspect cleanup.sh
pkmg inspect cleanup.sh --json
```

Run a skill:

```bash
pkmg run cleanup.sh
pkmg run cleanup.sh arg1 arg2
```

Launch the Web UI:

```bash
pkmg ui
```

## Web UI

The Web UI is intentionally lightweight.

- Backend: Go HTTP server + WebSocket command channel
- Frontend: single embedded HTML page
- UI libraries: loaded from CDN

Current local skill management flow:

- search local skills
- create a new skill from the editor modal
- edit existing skills
- copy skills with system-style copy suffixes
- switch back to previous versions
- open the containing folder in the OS file browser

The marketplace tab is intentionally hidden for now.

## Agent-Friendly Usage

If you are integrating `pkmg` into an agent workflow, prefer the structured CLI surface first.

Recommended commands:

```bash
pkmg list --json
pkmg search "<query>" --json
pkmg inspect "<relative-path>" --json
pkmg run "<relative-path>" [args...]
```

This gives an agent a predictable local capability inventory without requiring it to parse the UI.

For repo-specific agent guidance, see [AGENTS.md](./AGENTS.md).

## Project Layout

```text
cmd/              Cobra commands and backend logic
cmd/webui/        Embedded Web UI assets
data/scripts/     Managed local skill scripts
data/.pkmg/       Local metadata and version snapshots
platform/         Reserved for future distribution wrappers
```

## Testing

Build and smoke test locally:

```bash
go build ./...
pkmg init
pkmg list --json
pkmg ui
```

Useful manual checks:

- create a skill from the UI
- edit and save it
- copy it
- restore an earlier version
- verify `inspect --json` reflects the latest state

## Distribution Plan

The core product should stay a Go binary.

Recommended release order:

1. `go install`
2. GitHub Releases with multi-platform binaries
3. Homebrew
4. Windows package managers
5. optional thin npm wrapper for Node-first environments

The npm path, if added, should be a distribution wrapper around the Go binary, not a reimplementation of the core logic.

## Development

Build:

```bash
go build ./...
```

Run locally:

```bash
go run . --help
go run . ui
```

Format:

```bash
gofmt -w cmd/*.go main.go
```

## Status

`pkmg` is still early-stage. The current focus is:

- strong local skill management
- agent-friendly discovery and inspection
- lightweight UX instead of a heavy platform

Feedback and iteration ideas are welcome.
