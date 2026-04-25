#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
VERSION_INPUT="${1:-}"

if [[ -z "${VERSION_INPUT}" ]]; then
  echo "Usage: scripts/publish/go.sh <version>"
  echo "Example: scripts/publish/go.sh v0.2.1"
  exit 1
fi

if [[ "${VERSION_INPUT}" =~ ^v ]]; then
  TAG="${VERSION_INPUT}"
else
  TAG="v${VERSION_INPUT}"
fi

cd "${ROOT_DIR}"

if ! command -v gh >/dev/null 2>&1; then
  echo "Missing dependency: gh"
  exit 1
fi

if ! gh auth status >/dev/null 2>&1; then
  echo "GitHub CLI is not authenticated. Run: gh auth login"
  exit 1
fi

if ! git diff --quiet || ! git diff --cached --quiet; then
  echo "Working tree is not clean. Commit or stash changes before publishing."
  exit 1
fi

CURRENT_COMMIT="$(git rev-parse HEAD)"
if ! git rev-parse "${TAG}" >/dev/null 2>&1; then
  echo "Tag does not exist locally: ${TAG}"
  echo "Create it first, for example:"
  echo "  git tag -a ${TAG} -m \"${TAG}\""
  exit 1
fi

TAG_COMMIT="$(git rev-list -n 1 "${TAG}")"
if [[ "${CURRENT_COMMIT}" != "${TAG_COMMIT}" ]]; then
  echo "Tag ${TAG} does not point at HEAD."
  echo "HEAD: ${CURRENT_COMMIT}"
  echo "TAG : ${TAG_COMMIT}"
  exit 1
fi

if ! git ls-remote --exit-code --tags origin "refs/tags/${TAG}" >/dev/null 2>&1; then
  echo "Remote tag not found on origin: ${TAG}"
  echo "Push it first:"
  echo "  git push origin ${TAG}"
  exit 1
fi

export GITHUB_TOKEN="${GITHUB_TOKEN:-$(gh auth token)}"

echo "Publishing Go release for ${TAG}..."
go run github.com/goreleaser/goreleaser/v2@v2.15.4 release --clean

