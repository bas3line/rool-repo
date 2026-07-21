#!/bin/sh
# Verify the human portal and canonical agent-readable Sandbox documentation set.
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
DOCS=$ROOT/public/docs/sandbox
WATCHMAN_DOCS=$ROOT/public/docs/watchman

for file in theme.css tools.css tools.js docs/docs.css docs/docs.js
do
  [ -s "$ROOT/public/$file" ] || {
    printf '%s\n' "missing or empty shared frontend asset: $file" >&2
    exit 1
  }
done

for file in \
  index.html index.md llms.txt overview.md architecture.md aegis.md cli.md mcp.md \
  api.md agents.md configuration.md deployment.md operations.md security.md \
  runtime-driver.md development.md roadmap.md reporting.md tunnels.md \
  how-to-setup/server.md how-to-setup/client.md how-to-setup/custom-public-domains.md
do
  [ -s "$DOCS/$file" ] || {
    printf '%s\n' "missing or empty Sandbox documentation file: $file" >&2
    exit 1
  }
done

grep -F 'href="/docs/sandbox/index.md"' "$DOCS/index.html" >/dev/null
grep -F 'href="/docs/sandbox/llms.txt"' "$DOCS/index.html" >/dev/null
grep -F 'href="/theme.css?v=20260720-1"' "$DOCS/index.html" >/dev/null
grep -F 'href="/docs/docs.css?v=20260720-1"' "$DOCS/index.html" >/dev/null
grep -F 'src="/docs/docs.js?v=20260720-1"' "$DOCS/index.html" >/dev/null
grep -F 'https://tools.yshubham.com/docs/sandbox/mcp.md' "$DOCS/index.md" >/dev/null
grep -F 'https://tools.yshubham.com/docs/sandbox/how-to-setup/server.md' "$DOCS/index.md" >/dev/null
grep -F 'href="/docs/sandbox/how-to-setup/custom-public-domains.md"' "$DOCS/index.html" >/dev/null
grep -F '<strong>12 tools</strong>' "$DOCS/index.html" >/dev/null
grep -F 'https://tools.yshubham.com/docs/sandbox/security.md' "$DOCS/llms.txt" >/dev/null
grep -F '<li><a href="https://docs.yshubham.com/">Docs</a></li>' "$ROOT/public/index.html" >/dev/null
grep -F 'href="https://docs.yshubham.com/"' "$ROOT/public/index.html" >/dev/null
grep -F 'docs.yshubham.com ↗' "$ROOT/public/index.html" >/dev/null

[ -s "$WATCHMAN_DOCS/index.html" ] || {
  printf '%s\n' 'missing or empty Watchman documentation portal' >&2
  exit 1
}
[ -s "$WATCHMAN_DOCS/reference.md" ] || {
  printf '%s\n' 'missing or empty Watchman raw Markdown reference' >&2
  exit 1
}
grep -F 'href="/theme.css?v=20260720-1"' "$WATCHMAN_DOCS/index.html" >/dev/null
grep -F 'href="/docs/docs.css?v=20260720-1"' "$WATCHMAN_DOCS/index.html" >/dev/null
grep -F 'src="/docs/docs.js?v=20260720-1"' "$WATCHMAN_DOCS/index.html" >/dev/null
grep -F 'href="/docs/watchman/reference.md"' "$WATCHMAN_DOCS/index.html" >/dev/null

grep -F 'href="/theme.css?v=20260721-2"' "$ROOT/public/index.html" >/dev/null
grep -F 'href="/registry.css?v=20260721-4"' "$ROOT/public/index.html" >/dev/null
grep -F 'src="/tools.js?v=20260721-1"' "$ROOT/public/index.html" >/dev/null
grep -F 'Codex' "$ROOT/public/index.html" >/dev/null
grep -F 'Claude Code' "$ROOT/public/index.html" >/dev/null
grep -F 'OpenCode' "$ROOT/public/index.html" >/dev/null
grep -F 'Gemini CLI' "$ROOT/public/index.html" >/dev/null

printf '%s\n' 'Tools frontend and Sandbox/Watchman documentation bundle are complete'
