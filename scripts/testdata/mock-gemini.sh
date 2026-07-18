#!/bin/sh
set -eu

if [ "${1:-}" = mcp ] && [ "${2:-}" = list ]; then
  exit 0
fi
printf 'gemini' >> "$SANDBOX_TEST_LOG"
printf ' %s' "$@" >> "$SANDBOX_TEST_LOG"
printf '\n' >> "$SANDBOX_TEST_LOG"
