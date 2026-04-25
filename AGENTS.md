# AGENTS.md

This file explains how agents should interact with `pkmg`.

## What pkmg is

`pkmg` is a local-first skill manager for reusable scripts.

Its job is to give agents a predictable way to discover, inspect, and run local capabilities that are already installed on the machine, while keeping human editing and version management simple through a lightweight Web UI.

## Prefer the CLI JSON surface

When an agent needs capability discovery or execution planning, use the structured CLI first.

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
- it keeps local automation predictable
- it lets humans keep the UI as a management surface, not an integration surface

## Path model

Managed scripts do not need to live inside the repository.

Default directory model:

- config root: `os.UserConfigDir()/pkmg`
- data root: `PKMG_DATA_DIR`, configured `dataPath`, or `os.UserConfigDir()/pkmg`
- scripts root: `<data-root>/scripts`
- version root: `<data-root>/.pkmg`

Script identities are always relative to the scripts root.

Examples:

- `cleanup.sh`
- `team/deploy.sh`

When invoking `pkmg inspect` or `pkmg run`, use those relative paths.

## Discovery workflow

Recommended flow for agents:

1. Call `pkmg list --json` to get the local inventory.
2. Narrow candidates with `pkmg search "<query>" --json` when needed.
3. Inspect the chosen skill with `pkmg inspect "<path>" --json`.
4. Only then run it with `pkmg run "<path>" [args...]`.

This reduces accidental execution and gives the agent a chance to verify metadata first.

## When to use the Web UI

Use `pkmg ui` when a human needs interactive management such as:

- creating a new skill
- editing an existing skill in a modal editor
- copying a skill
- restoring an older version
- opening the skill directory in the OS file browser

Agents should still prefer the CLI JSON surface for automation.

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
- The Web UI uses a local WebSocket schema, but that schema is currently an internal UI transport, not a public contract.
- Managed skill creation is currently exposed through the UI, not through a dedicated public CLI create command.

## Safe assumptions

Agents may safely assume:

- `pkmg list --json` returns the current local inventory
- `pkmg search --json` returns filtered candidates
- `pkmg inspect --json` returns metadata plus content details
- `pkmg run` executes a local script via shell
- `pkmg ui` starts a local-only management interface

Agents should not assume:

- cloud sync
- remote registry support
- public marketplace availability
- a stable third-party integration protocol yet
