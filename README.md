# <div align="center">pkmg</div>

<div align="center">

_Local-first skill manager for humans and coding agents._

[![release](https://img.shields.io/github/v/release/CoderSerio/pokemand-go?display_name=tag&style=flat-square)](https://github.com/CoderSerio/pokemand-go/releases)
[![ci](https://img.shields.io/github/actions/workflow/status/CoderSerio/pokemand-go/ci.yml?branch=main&style=flat-square&label=ci)](https://github.com/CoderSerio/pokemand-go/actions/workflows/ci.yml)
[![go](https://img.shields.io/github/go-mod/go-version/CoderSerio/pokemand-go?style=flat-square)](https://go.dev/)
[![license](https://img.shields.io/github/license/CoderSerio/pokemand-go?style=flat-square)](./LICENSE)
[![local-first](https://img.shields.io/badge/local--first-yes-0f766e?style=flat-square)](#why-pkmg)
[![agent-ready](https://img.shields.io/badge/agent-ready-yes-1d4ed8?style=flat-square)](./AGENTS.md)

[简体中文](./README_ZH.md) · [Quick Start](#quick-start) · [Web UI](#web-ui) · [Agent Usage](#agent-friendly-usage) · [Testing](#testing)

</div>

`pkmg` turns scattered local scripts into a reusable local skill inventory.

It is a lightweight script-layer manager for the local-agent workflow: discover, inspect, edit, version, and run existing skills without turning the stack into a heavy platform.

`pkmg` is the CLI face of `pokemand-go`, short for a Go-built "pocket command" manager.

## Why pkmg

`skill` is becoming a more natural local pattern for agent workflows.

That helps with locality, but it still leaves one annoying layer unresolved:

- scripts are duplicated across skills and projects
- local capabilities become hard to browse and harder to trust
- once version history matters, ad-hoc files become messy fast
- agents need structured discovery, not "please read this folder"

`pkmg` sits exactly on that layer.

It is not trying to replace agents, skill files, or project-specific automation. It is a lightweight manager for the script side of local skills: one place to organize them, version them, search them, inspect them, and reuse them across workflows.

## What It Feels Like

- Local-first by default. No registry or cloud dependency is required.
- Agent-friendly from day one. JSON commands expose a stable discovery surface.
- Light enough to stay out of the way. Go backend, embedded page, CDN UI dependencies.
- Version-aware. Script edits automatically produce local snapshots.
- Human-usable. A tiny Web UI gives you create, edit, copy, restore, and folder-open flows.

## Current Features

- Initialize a user-level skill workspace
- Store managed scripts under a dedicated local data directory
- List, search, inspect, and run skills from the CLI
- Emit structured JSON for agent workflows
- Launch a lightweight local Web UI with:
  - local skill listing
  - search
  - create
  - edit
  - copy
  - version switching
  - open containing folder

## Installation

### Go install

```bash
go install github.com/CoderSerio/pokemand-go@latest
```

### Build from source

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

## Default Directory Model

By default, `pkmg` uses user-level directories instead of repo-local `data/`.

- Config root: `os.UserConfigDir()/pkmg`
- Data root: `PKMG_DATA_DIR`, or configured `dataPath`, or `os.UserConfigDir()/pkmg`
- Scripts root: `<data-root>/scripts`
- Version snapshots: `<data-root>/.pkmg`

Environment overrides:

```bash
export PKMG_CONFIG_DIR=/your/custom/config
export PKMG_DATA_DIR=/your/custom/data
```

This keeps managed skills outside the repository by default, which is a better fit for reusable local tooling.

## Web UI

- Backend: Go HTTP server + WebSocket command transport
- Frontend: a single embedded page
- UI dependencies: loaded from CDN to keep the binary lean

Current local skill workflow:

- search local skills
- create a new skill from the editor modal
- edit an existing skill with a lightweight code view
- copy a skill with system-style copy naming
- restore an older version
- open the containing folder in the OS file browser

## Agent-Friendly Usage

If you are integrating `pkmg` into an agent workflow, prefer the structured CLI surface first.

Recommended commands:

```bash
pkmg list --json
pkmg search "<query>" --json
pkmg inspect "<relative-path>" --json
pkmg run "<relative-path>" [args...]
```

That gives an agent a predictable local capability inventory without scraping the UI.

For repo-specific guidance aimed at agents, see [AGENTS.md](./AGENTS.md).

## Testing

Automated:

```bash
go test ./...
```

Quick smoke flow:

```bash
go build ./...
pkmg init
pkmg list --json
pkmg ui
```
