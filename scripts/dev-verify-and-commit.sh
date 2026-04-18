#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

SKIP_TESTS=0
if [[ "${1:-}" == "--skip-tests" ]]; then
  SKIP_TESTS=1
fi

run_backend_tests() {
  echo "== Backend: go test ./... -count=1 -race -timeout=2m =="
  ( cd "$ROOT/backend" && CGO_ENABLED=1 go test ./... -count=1 -race -timeout=2m )
}

run_frontend_tests() {
  echo "== Frontend: pnpm test -- --passWithNoTests =="
  ( cd "$ROOT/frontend" && pnpm test -- --passWithNoTests )
}

if [[ "$SKIP_TESTS" -eq 0 ]]; then
  run_backend_tests
  run_frontend_tests
fi

if [[ -z "${GONS_COMMIT_SUBJECT:-}" ]]; then
  echo "ERROR: set GONS_COMMIT_SUBJECT (commit subject)" >&2
  exit 1
fi
if [[ -z "${GONS_COMMIT_EXPLANATION:-}" ]]; then
  echo "ERROR: set GONS_COMMIT_EXPLANATION (commit body explanation)" >&2
  exit 1
fi

WHEN="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
STATUS="$(git -C "$ROOT" status -sb)"

if git -C "$ROOT" diff --cached --quiet; then
  CACHED_DISPLAY='_None._'
else
  CACHED_DISPLAY="$(git -C "$ROOT" diff --cached --stat)"
fi

if git -C "$ROOT" diff --quiet; then
  UNSTAGED_DISPLAY='_None._'
else
  UNSTAGED_DISPLAY="$(git -C "$ROOT" diff --stat)"
fi

MSG_FILE="$(mktemp)"
trap 'rm -f "$MSG_FILE"' EXIT

{
  printf '%s\n\n' "$GONS_COMMIT_SUBJECT"
  printf '%s\n\n' "$GONS_COMMIT_EXPLANATION"
  printf '## Verification\n'
  printf -- '- When (UTC): %s\n' "$WHEN"
  printf '%s\n' '- Backend: go test ./... -count=1 -race -timeout=2m (passed)'
  printf '%s\n\n' '- Frontend: pnpm test -- --passWithNoTests (passed)'
  printf '## Git analysis\n'
  printf '### Branch status\n```\n%s\n```\n\n' "$STATUS"
  printf '%s\n' '### Staged changes (included in this commit)'
  if [[ "$CACHED_DISPLAY" == '_None._' ]]; then
    printf '%s\n\n' '_None._'
  else
    printf '```\n%s\n```\n\n' "$CACHED_DISPLAY"
  fi
  printf '%s\n' '### Unstaged changes (not in this commit)'
  if [[ "$UNSTAGED_DISPLAY" == '_None._' ]]; then
    printf '%s\n' '_None._'
  else
    printf '```\n%s\n```\n' "$UNSTAGED_DISPLAY"
  fi
} >"$MSG_FILE"

if git -C "$ROOT" diff --cached --quiet; then
  if [[ "${GONS_ALLOW_EMPTY_COMMIT:-}" == "1" ]]; then
    echo "WARN: no staged changes; empty commit (GONS_ALLOW_EMPTY_COMMIT=1)" >&2
    git -C "$ROOT" commit --allow-empty -F "$MSG_FILE"
  else
    echo "ERROR: No staged changes. git add files first, or set GONS_ALLOW_EMPTY_COMMIT=1" >&2
    exit 1
  fi
else
  git -C "$ROOT" commit -F "$MSG_FILE"
fi

if [[ "${GONS_GH_PR:-}" == "1" ]]; then
  gh pr comment --body-file "$MSG_FILE" || echo "WARN: gh pr comment failed" >&2
fi

if [[ "${GONS_SKIP_PUSH:-}" == "1" ]]; then
  echo "GONS_SKIP_PUSH=1 — skipping git push"
  exit 0
fi

echo "== git push =="
git -C "$ROOT" push
