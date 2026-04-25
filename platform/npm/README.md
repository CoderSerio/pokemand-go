# @seriohhhh/pkmg

This package is a thin npm wrapper around the `pkmg` Go binary.

It downloads the matching GitHub Release asset for the current package version during `postinstall`, then exposes a `pkmg` executable on your `PATH`.

## Install

```bash
npm install -g @seriohhhh/pkmg
```

## Notes

- This package does not reimplement `pkmg` in JavaScript.
- The actual runtime is the Go binary published from the main repository.
- This wrapper currently follows the GitHub Release assets produced by `pokemand-go`.

Main project:

- https://github.com/CoderSerio/pokemand-go
