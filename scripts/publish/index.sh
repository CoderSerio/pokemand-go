#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
VERSION_INPUT="${1:-}"
TARGET="${2:-go}"

if [[ -z "${VERSION_INPUT}" ]]; then
  echo "Usage: scripts/publish/index.sh <version> [go]"
  echo "Example: scripts/publish/index.sh v0.2.1"
  exit 1
fi

case "${TARGET}" in
  go)
    "${ROOT_DIR}/scripts/publish/go.sh" "${VERSION_INPUT}"
    ;;
  *)
    echo "Unknown target: ${TARGET}"
    echo "Expected: go"
    exit 1
    ;;
esac
