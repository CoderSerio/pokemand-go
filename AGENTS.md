# AGENTS.md

This file explains how agents should interact with `pkmg`.

## What pkmg is

`pkmg` is a local-first skill manager built around scripts stored under `data/scripts/`.

It is useful when an agent needs to:

- discover what local reusable skills already exist
- inspect a skill before execution
- run a known local skill
- ask a user to open the Web UI for interactive management

## Prefer the CLI JSON surface

When an agent needs local capability discovery, use the structured CLI first.

Preferred commands:

```bash
pkmg list --json
pkmg search "<query>" --json
pkmg inspect "<relative-path>" --json
pkmg run "<relative-path>" [args...]
```

Why:

- the output is machine-readable
- it avoids scraping the UI
- it keeps the local workflow predictable

## Expected path model

- Managed scripts live under `data/scripts/`
- Script identities are relative paths from that directory
- Example:
  - `cleanup.sh`
  - `team/deploy.sh`

When invoking `pkmg inspect` or `pkmg run`, use those relative paths.

## Discovery workflow

Recommended flow for agents:

1. Call `pkmg list --json` to get the current local inventory.
2. If needed, narrow candidates with `pkmg search "<query>" --json`.
3. Inspect a chosen skill with `pkmg inspect "<path>" --json`.
4. Only then run it with `pkmg run "<path>" [args...]`.

## When to use the Web UI

Use `pkmg ui` when a human needs interactive management such as:

- creating a new skill
- editing an existing skill in a modal editor
- copying a skill
- restoring an older version
- opening the skill directory in the OS file browser

Agents should still prefer the CLI JSON surface for non-interactive automation.

## Current command surface

Useful commands today:

```bash
pkmg init
pkmg open <path>
pkmg list [--json] [--search <query>]
pkmg search <query> [--json]
pkmg inspect <path> [--json]
pkmg run <path> [args...]
pkmg ui
```

## Important limitations

- There is no stable remote API yet.
- The marketplace tab is intentionally hidden.
- The Web UI uses a local WebSocket schema, but that schema is currently an internal UI transport, not a public remote integration contract.
- Managed skill creation is currently exposed through the UI, not through a dedicated public CLI subcommand.

## Safe assumptions

Agents may safely assume:

- `pkmg list --json` returns the current known inventory
- `pkmg inspect --json` returns metadata plus a preview or content details
- `pkmg run` executes a local script via shell
- `pkmg ui` starts a local-only management interface

Agents should not assume:

- cloud sync
- package registry support
- public marketplace availability
- stable third-party distribution wrappers yet
