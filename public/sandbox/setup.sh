#!/bin/sh
# Full Sandbox workstation setup: binaries, agent skill, and Codex MCP registration.
set -eu
umask 077

fail() {
  printf '%s\n' "sandbox setup: $*" >&2
  exit 1
}

BASE_URL=${SANDBOX_INSTALL_BASE_URL:-https://tools.yshubham.com}
INSTALL_SKILL=${SANDBOX_SETUP_SKILL:-1}
CONFIGURE_CODEX=${SANDBOX_SETUP_CODEX_MCP:-1}

usage() {
  printf '%s\n' 'usage: setup.sh [--no-skill] [--no-codex-mcp]'
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --no-skill) INSTALL_SKILL=0 ;;
    --no-codex-mcp) CONFIGURE_CODEX=0 ;;
    -h|--help) usage; exit 0 ;;
    *) usage >&2; fail "unknown option: $1" ;;
  esac
  shift
done

case "$INSTALL_SKILL" in 0|1) ;; *) fail "SANDBOX_SETUP_SKILL must be 0 or 1" ;; esac
case "$CONFIGURE_CODEX" in 0|1) ;; *) fail "SANDBOX_SETUP_CODEX_MCP must be 0 or 1" ;; esac
case "$BASE_URL" in https://*) ;; *) fail "SANDBOX_INSTALL_BASE_URL must use HTTPS" ;; esac
command -v curl >/dev/null 2>&1 || fail "missing required command: curl"
command -v mktemp >/dev/null 2>&1 || fail "missing required command: mktemp"
if [ "$INSTALL_SKILL" = 1 ]; then
  command -v npx >/dev/null 2>&1 || fail "npx is required for the agent skill; install Node.js or rerun with --no-skill"
fi

temporary=$(mktemp -d "${TMPDIR:-/tmp}/sandbox-setup.XXXXXX")
trap 'rm -rf "$temporary"' EXIT INT TERM
curl --fail --location --proto '=https' --proto-redir '=https' --silent --show-error \
  --output "$temporary/install.sh" -- "${BASE_URL%/}/sandbox/install.sh"
/bin/sh "$temporary/install.sh" </dev/null

if [ -n "${SANDBOX_INSTALL_DIR:-}" ]; then
  install_dir=$SANDBOX_INSTALL_DIR
elif [ -d /usr/local/bin ] && [ -w /usr/local/bin ]; then
  install_dir=/usr/local/bin
else
  [ -n "${HOME:-}" ] || fail "HOME is unset; set SANDBOX_INSTALL_DIR"
  install_dir=${XDG_BIN_HOME:-$HOME/.local/bin}
fi

if [ "$INSTALL_SKILL" = 1 ]; then
  printf '%s\n' "Installing the Sandbox skill for supported coding agents"
  npx --yes skills add bas3line/rool-repo --skill sandbox-platform --agent '*' --global --yes </dev/null
fi

if [ "$CONFIGURE_CODEX" = 1 ] && command -v codex >/dev/null 2>&1; then
  if codex mcp get sandbox </dev/null >/dev/null 2>&1; then
    printf '%s\n' "Codex MCP entry 'sandbox' already exists; leaving it unchanged"
  else
    codex mcp add sandbox -- "$install_dir/sandbox-mcp" </dev/null
    printf '%s\n' "Registered sandbox-mcp with Codex"
  fi
elif [ "$CONFIGURE_CODEX" = 1 ]; then
  printf '%s\n' "Codex is not installed; skipped Codex MCP registration"
fi

printf '%s\n' "Sandbox workstation setup is complete."
printf '%s\n' "Set SANDBOX_URL and SANDBOX_TOKEN in the environment that launches your agent."
printf '%s\n' "Verify: $install_dir/sandbox doctor"
