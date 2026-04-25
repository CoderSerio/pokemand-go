#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
VERSION_INPUT="${1:-}"
TARGET="${2:-all}"

if [[ -z "${VERSION_INPUT}" ]]; then
  echo "Usage: scripts/publish/index.sh <version> [go|npm|all]"
  echo "Example: scripts/publish/index.sh v0.2.1 all"
  exit 1
fi

case "${TARGET}" in
  go)
    "${ROOT_DIR}/scripts/publish/go.sh" "${VERSION_INPUT}"
    ;;
  npm)
    "${ROOT_DIR}/scripts/publish/npm.sh" "${VERSION_INPUT}"
    ;;
  all)
    "${ROOT_DIR}/scripts/publish/go.sh" "${VERSION_INPUT}"
    "${ROOT_DIR}/scripts/publish/npm.sh" "${VERSION_INPUT}"
    ;;
  *)
    echo "Unknown target: ${TARGET}"
    echo "Expected one of: go, npm, all"
    exit 1
    ;;
esac

