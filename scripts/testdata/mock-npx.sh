#!/bin/sh
set -eu
printf 'npx' >> "$SANDBOX_TEST_LOG"
printf ' %s' "$@" >> "$SANDBOX_TEST_LOG"
printf '\n' >> "$SANDBOX_TEST_LOG"
