#!/bin/sh
# Full Sandbox workstation setup: binaries, cross-agent skill, and detected MCP CLIs.
set -eu
umask 077

fail() {
  printf '%s\n' "sandbox setup: $*" >&2
  exit 1
}

BASE_URL=${SANDBOX_INSTALL_BASE_URL:-https://tools.yshubham.com}
INSTALL_SKILL=${SANDBOX_SETUP_SKILL:-1}
CONFIGURE_MCP=${SANDBOX_SETUP_MCP:-${SANDBOX_SETUP_CODEX_MCP:-1}}

usage() {
  printf '%s\n' 'usage: setup.sh [--no-skill] [--no-mcp]'
  printf '%s\n' '       --no-codex-mcp remains as a deprecated alias for --no-mcp'
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --no-skill) INSTALL_SKILL=0 ;;
    --no-mcp|--no-codex-mcp) CONFIGURE_MCP=0 ;;
    -h|--help) usage; exit 0 ;;
    *) usage >&2; fail "unknown option: $1" ;;
  esac
  shift
done

case "$INSTALL_SKILL" in 0|1) ;; *) fail "SANDBOX_SETUP_SKILL must be 0 or 1" ;; esac
case "$CONFIGURE_MCP" in 0|1) ;; *) fail "SANDBOX_SETUP_MCP must be 0 or 1" ;; esac
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

if [ "$CONFIGURE_MCP" = 1 ]; then
  detected_mcp_clients=0

  if command -v codex >/dev/null 2>&1; then
    detected_mcp_clients=$((detected_mcp_clients + 1))
    if codex mcp get sandbox </dev/null >/dev/null 2>&1; then
      printf '%s\n' "Codex MCP entry 'sandbox' already exists; leaving it unchanged"
    else
      codex mcp add sandbox -- "$install_dir/sandbox-mcp" </dev/null
      printf '%s\n' "Registered sandbox-mcp with Codex"
    fi
  fi

  if command -v claude >/dev/null 2>&1; then
    detected_mcp_clients=$((detected_mcp_clients + 1))
    if claude mcp get sandbox </dev/null >/dev/null 2>&1; then
      printf '%s\n' "Claude Code MCP entry 'sandbox' already exists; leaving it unchanged"
    else
      claude mcp add --scope user --transport stdio sandbox -- "$install_dir/sandbox-mcp" </dev/null
      printf '%s\n' "Registered sandbox-mcp with Claude Code"
    fi
  fi

  if command -v gemini >/dev/null 2>&1; then
    detected_mcp_clients=$((detected_mcp_clients + 1))
    gemini_mcp_list=$(gemini mcp list </dev/null 2>/dev/null || true)
    case "$gemini_mcp_list" in
      *sandbox*) printf '%s\n' "Gemini CLI MCP entry 'sandbox' already exists; leaving it unchanged" ;;
      *)
        gemini mcp add sandbox "$install_dir/sandbox-mcp" --scope user </dev/null
        printf '%s\n' "Registered sandbox-mcp with Gemini CLI"
        ;;
    esac
  fi

  if [ "$detected_mcp_clients" -eq 0 ]; then
    printf '%s\n' "No supported non-interactive MCP CLI detected; use ${BASE_URL%/}/sandbox/clients/index.md"
  else
    printf '%s\n' "Cursor, OpenCode, VS Code, Windsurf, Cline, Roo Code, Goose, and generic client templates:"
    printf '%s\n' "${BASE_URL%/}/sandbox/clients/index.md"
  fi
fi

printf '%s\n' "Sandbox workstation setup is complete."
printf '%s\n' "Set SANDBOX_URL and SANDBOX_TOKEN in the environment that launches your agent."
printf '%s\n' "Verify: $install_dir/sandbox doctor"
