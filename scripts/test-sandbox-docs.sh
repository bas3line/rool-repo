#!/bin/sh
# Verify the human portal and canonical agent-readable Sandbox documentation set.
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
DOCS=$ROOT/public/docs/sandbox

for file in \
  index.html index.md llms.txt overview.md architecture.md aegis.md cli.md mcp.md \
  api.md agents.md configuration.md deployment.md operations.md security.md \
  runtime-driver.md development.md roadmap.md reporting.md
do
  [ -s "$DOCS/$file" ] || {
    printf '%s\n' "missing or empty Sandbox documentation file: $file" >&2
    exit 1
  }
done

grep -F 'href="/docs/sandbox/index.md"' "$DOCS/index.html" >/dev/null
grep -F 'href="/docs/sandbox/llms.txt"' "$DOCS/index.html" >/dev/null
grep -F 'https://tools.yshubham.com/docs/sandbox/mcp.md' "$DOCS/index.md" >/dev/null
grep -F 'https://tools.yshubham.com/docs/sandbox/security.md' "$DOCS/llms.txt" >/dev/null
grep -F 'href="/docs/sandbox/"' "$ROOT/public/index.html" >/dev/null

printf '%s\n' 'Sandbox human and agent documentation bundle is complete'
