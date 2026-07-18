#!/bin/sh
# Exercise the public installer without network access using the release for this host.
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
  printf '%s\n' "installer test skipped: $archive is not published" >&2
  exit 0
}

temporary=$(mktemp -d "${TMPDIR:-/tmp}/sandbox-installer-test.XXXXXX")
trap 'rm -rf "$temporary"' EXIT INT TERM
mkdir -p "$temporary/bin" "$temporary/install"
cp "$ROOT/scripts/testdata/mock-curl.sh" "$temporary/bin/curl"
chmod 0755 "$temporary/bin/curl"

export SANDBOX_TEST_VERSION=$VERSION
export SANDBOX_TEST_FIXTURE=$fixture
PATH="$temporary/bin:/usr/bin:/bin:/usr/sbin:/sbin" \
SANDBOX_INSTALL_BASE_URL=https://tools.test \
SANDBOX_INSTALL_DIR="$temporary/install" \
/bin/sh "$ROOT/public/sandbox/install.sh"

"$temporary/install/sandbox" --version
"$temporary/install/sandbox-mcp" --version
"$temporary/install/sandboxd" --version
printf '%s\n' "Sandbox installer integration test passed for $os/$arch"
