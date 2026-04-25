#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
VERSION_INPUT="${1:-}"
NPM_DIR="${ROOT_DIR}/platform/npm"

if [[ -z "${VERSION_INPUT}" ]]; then
  echo "Usage: scripts/publish/npm.sh <version>"
  echo "Example: scripts/publish/npm.sh v0.2.1"
  exit 1
fi

if [[ "${VERSION_INPUT}" =~ ^v ]]; then
  VERSION="${VERSION_INPUT#v}"
else
  VERSION="${VERSION_INPUT}"
fi

cd "${ROOT_DIR}"

if [[ ! -d "${NPM_DIR}" ]]; then
  echo "Missing npm wrapper directory: ${NPM_DIR}"
  echo "Create platform/npm before publishing to npm."
  exit 1
fi

if [[ ! -f "${NPM_DIR}/package.json" ]]; then
  echo "Missing package.json: ${NPM_DIR}/package.json"
  exit 1
fi

if ! command -v npm >/dev/null 2>&1; then
  echo "Missing dependency: npm"
  exit 1
fi

NPM_WHOAMI_OUTPUT=""
if ! NPM_WHOAMI_OUTPUT="$(npm whoami 2>&1)"; then
  printf '%s\n' "${NPM_WHOAMI_OUTPUT}"
  echo "npm is not authenticated. Run: npm adduser"
  exit 1
fi

if ! git diff --quiet || ! git diff --cached --quiet; then
  echo "Working tree is not clean. Commit or stash changes before publishing."
  exit 1
fi

cd "${NPM_DIR}"

CURRENT_VERSION="$(node -p "require('./package.json').version")"
if [[ "${CURRENT_VERSION}" != "${VERSION}" ]]; then
  echo "package.json version mismatch."
  echo "Expected: ${VERSION}"
  echo "Current : ${CURRENT_VERSION}"
  echo "Update the npm wrapper version before publishing."
  exit 1
fi

echo "Publishing npm package version ${VERSION} from ${NPM_DIR}..."
npm publish --access public
