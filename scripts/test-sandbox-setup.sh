#!/bin/sh
# Exercise the complete setup flow with local release fixtures and fake agent CLIs.
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
VERSION=v0.1.0
os=$(uname -s | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)
case "$arch" in
  x86_64|amd64) arch=amd64 ;;
  aarch64|arm64) arch=arm64 ;;
  *) printf '%s\n' "unsupported test architecture: $arch" >&2; exit 1 ;;
esac
fixture="$ROOT/public/sandbox/releases/$VERSION"
archive="sandbox_${VERSION}_${os}_${arch}.tar.gz"
[ -f "$fixture/$archive" ] || {
  printf '%s\n' "setup test skipped: $archive is not published" >&2
  exit 0
}

temporary=$(mktemp -d "${TMPDIR:-/tmp}/sandbox-setup-test.XXXXXX")
trap 'rm -rf "$temporary"' EXIT INT TERM
mkdir -p "$temporary/bin" "$temporary/install"
cp "$ROOT/scripts/testdata/mock-curl.sh" "$temporary/bin/curl"
cp "$ROOT/scripts/testdata/mock-npx.sh" "$temporary/bin/npx"
cp "$ROOT/scripts/testdata/mock-codex.sh" "$temporary/bin/codex"
chmod 0755 "$temporary/bin/curl" "$temporary/bin/npx" "$temporary/bin/codex"
: > "$temporary/actions.log"

export SANDBOX_TEST_VERSION=$VERSION
export SANDBOX_TEST_FIXTURE=$fixture
export SANDBOX_TEST_INSTALLER=$ROOT/public/sandbox/install.sh
export SANDBOX_TEST_LOG=$temporary/actions.log
PATH="$temporary/bin:/usr/bin:/bin:/usr/sbin:/sbin" \
SANDBOX_INSTALL_BASE_URL=https://tools.test \
SANDBOX_INSTALL_DIR="$temporary/install" \
/bin/sh -s -- < "$ROOT/public/sandbox/setup.sh"

"$temporary/install/sandbox" --version
"$temporary/install/sandbox-mcp" --version
grep -F "npx --yes skills add bas3line/rool-repo --skill sandbox-platform --agent * --global --yes" "$temporary/actions.log" >/dev/null
grep -F "codex mcp add sandbox -- $temporary/install/sandbox-mcp" "$temporary/actions.log" >/dev/null
printf '%s\n' "Sandbox full setup integration test passed for $os/$arch"
