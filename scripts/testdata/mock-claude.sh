#!/bin/sh
set -eu

if [ "${1:-}" = mcp ] && [ "${2:-}" = get ]; then
  exit 1
fi
printf 'claude' >> "$SANDBOX_TEST_LOG"
printf ' %s' "$@" >> "$SANDBOX_TEST_LOG"
printf '\n' >> "$SANDBOX_TEST_LOG"
