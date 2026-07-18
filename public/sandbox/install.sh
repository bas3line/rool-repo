#!/bin/sh
# Canonical Sandbox binary installer: https://tools.yshubham.com/sandbox/install.sh
set -eu
umask 077

fail() {
  printf '%s\n' "sandbox installer: $*" >&2
  exit 1
}

BASE_URL=${SANDBOX_INSTALL_BASE_URL:-https://tools.yshubham.com}
VERSION=${SANDBOX_VERSION:-latest}
BASE_URL=${BASE_URL%/}

case "$BASE_URL" in
  https://*) ;;
  *) fail "SANDBOX_INSTALL_BASE_URL must use HTTPS" ;;
esac

for required_command in awk curl install mkdir mktemp rm tar tr uname; do
  command -v "$required_command" >/dev/null 2>&1 || fail "missing required command: $required_command"
done
if ! command -v sha256sum >/dev/null 2>&1 && ! command -v shasum >/dev/null 2>&1; then
  fail "missing required checksum command: sha256sum or shasum"
fi

if [ "$VERSION" = "latest" ]; then
  VERSION=$(curl --fail --location --proto '=https' --proto-redir '=https' --silent --show-error -- "$BASE_URL/sandbox/latest") \
    || fail "cannot resolve the latest Sandbox version"
fi
case "$VERSION" in
  v[0-9]*) ;;
  *) fail "SANDBOX_VERSION must begin with v and a digit" ;;
esac
case "$VERSION" in
  *[!A-Za-z0-9._-]*) fail "SANDBOX_VERSION contains unsafe characters" ;;
esac

if [ -n "${SANDBOX_INSTALL_DIR:-}" ]; then
  INSTALL_DIR=$SANDBOX_INSTALL_DIR
elif [ -d /usr/local/bin ] && [ -w /usr/local/bin ]; then
  INSTALL_DIR=/usr/local/bin
else
  [ -n "${HOME:-}" ] || fail "HOME is unset; set SANDBOX_INSTALL_DIR"
  INSTALL_DIR=${XDG_BIN_HOME:-$HOME/.local/bin}
fi
case "$INSTALL_DIR" in
  /*) ;;
  *) fail "SANDBOX_INSTALL_DIR must be an absolute path" ;;
esac

os=$(uname -s | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)
case "$arch" in
  x86_64|amd64) arch=amd64 ;;
  aarch64|arm64) arch=arm64 ;;
  *) fail "unsupported architecture: $arch" ;;
esac
case "$os" in
  linux|darwin) ;;
  *) fail "unsupported operating system: $os" ;;
esac

archive="sandbox_${VERSION}_${os}_${arch}.tar.gz"
url="$BASE_URL/sandbox/releases/$VERSION/$archive"
temporary=$(mktemp -d "${TMPDIR:-/tmp}/sandbox-install.XXXXXX")
trap 'rm -rf "$temporary"' EXIT INT TERM

printf '%s\n' "Installing Sandbox $VERSION for $os/$arch"
curl --fail --location --proto '=https' --proto-redir '=https' --silent --show-error \
  --output "$temporary/$archive" -- "$url?v=$VERSION"
curl --fail --location --proto '=https' --proto-redir '=https' --silent --show-error \
  --output "$temporary/$archive.sha256" -- "$url.sha256?v=$VERSION" \
  || fail "checksum file is unavailable"

expected=$(awk '
  NR == 1 { print $1; next }
  { exit 1 }
  END { if (NR != 1) exit 1 }
' "$temporary/$archive.sha256") || fail "malformed checksum file"
[ "${#expected}" -eq 64 ] || fail "malformed SHA-256 digest"
case "$expected" in
  *[!0-9A-Fa-f]*) fail "malformed SHA-256 digest" ;;
esac
if command -v sha256sum >/dev/null 2>&1; then
  actual=$(sha256sum "$temporary/$archive" | awk '{print $1}')
else
  actual=$(shasum -a 256 "$temporary/$archive" | awk '{print $1}')
fi
[ "$expected" = "$actual" ] || fail "checksum mismatch"

member_names=$(tar -tzf "$temporary/$archive") || fail "cannot inspect release archive"
expected_members=$(printf '%s\n' sandbox sandboxd sandbox-mcp)
[ "$member_names" = "$expected_members" ] || fail "archive must contain only sandbox, sandboxd, and sandbox-mcp"
member_listing=$(tar -tvzf "$temporary/$archive") || fail "cannot inspect release archive"
printf '%s\n' "$member_listing" | awk '
  substr($1, 1, 1) == "-" { regular += 1; next }
  { bad = 1 }
  END { if (NR != 3 || regular != 3 || bad) exit 1 }
' >/dev/null || fail "archive members must be exactly three regular files"

tar -xzf "$temporary/$archive" -C "$temporary" sandbox sandboxd sandbox-mcp
mkdir -p "$INSTALL_DIR"
for binary in sandbox sandboxd sandbox-mcp; do
  [ -f "$temporary/$binary" ] || fail "archive is missing $binary"
  [ ! -L "$temporary/$binary" ] || fail "archive contains a symbolic-link $binary"
  install -m 0755 "$temporary/$binary" "$INSTALL_DIR/$binary"
done

printf '%s\n' "Installed sandbox, sandboxd, and sandbox-mcp $VERSION into $INSTALL_DIR"
case ":${PATH:-}:" in
  *:"$INSTALL_DIR":*) ;;
  *) printf '%s\n' "Add $INSTALL_DIR to PATH before running sandbox." ;;
esac
printf '%s\n' "Run: $INSTALL_DIR/sandbox --help"
