#!/bin/sh
# Watchman release installer served at https://tools.yshubham.com/watchman/install.sh.
set -eu
umask 077

fail() {
  printf '%s\n' "watchman: $*" >&2
  exit 1
}

# WATCHMAN_* is the public contract. Keep GPU_WATCHMAN_* as a migration alias.
BASE_URL=${WATCHMAN_BASE_URL:-${GPU_WATCHMAN_BASE_URL:-https://tools.yshubham.com}}
VERSION=${WATCHMAN_VERSION:-${GPU_WATCHMAN_VERSION:-v0.8.3}}
INSTALL_DIR=${WATCHMAN_INSTALL_DIR:-${GPU_WATCHMAN_INSTALL_DIR:-/usr/local/bin}}
BASE_URL=${BASE_URL%/}

case "$BASE_URL" in
  https://*) ;;
  *) fail "WATCHMAN_BASE_URL must use HTTPS" ;;
esac
case "$VERSION" in
  v[0-9]*) ;;
  *) fail "WATCHMAN_VERSION must begin with v and a digit" ;;
esac
case "$VERSION" in
  *[!A-Za-z0-9._-]*) fail "WATCHMAN_VERSION contains unsafe characters" ;;
esac
case "$INSTALL_DIR" in
  /*) ;;
  *) fail "WATCHMAN_INSTALL_DIR must be an absolute path" ;;
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

for required_command in awk curl install tar; do
  command -v "$required_command" >/dev/null 2>&1 || fail "missing required command: $required_command"
done
if ! command -v sha256sum >/dev/null 2>&1 && ! command -v shasum >/dev/null 2>&1; then
  fail "missing required checksum command: sha256sum or shasum"
fi

archive="gpu-watchman_${os}_${arch}.tar.gz"
url="$BASE_URL/watchman/releases/$VERSION/$archive"
cache_key="?v=$VERSION"
tmp=$(mktemp -d "${TMPDIR:-/tmp}/watchman.XXXXXX")
trap 'rm -rf "$tmp"' EXIT INT TERM

printf '%s\n' "Installing Watchman $VERSION for $os/$arch"
curl --fail --location --proto '=https' --proto-redir '=https' --silent --show-error \
  --output "$tmp/$archive" -- "$url$cache_key"
curl --fail --location --proto '=https' --proto-redir '=https' --silent --show-error \
  --output "$tmp/$archive.sha256" -- "$url.sha256$cache_key" || fail "checksum file is unavailable"

expected=$(awk -v name="$archive" '
  NR == 1 {
    file = $2
    sub(/^\*/, "", file)
    if (NF != 2 || file != name) exit 1
    print $1
    next
  }
  { exit 1 }
  END { if (NR != 1) exit 1 }
' "$tmp/$archive.sha256") || fail "malformed checksum file"
[ "${#expected}" -eq 64 ] || fail "malformed SHA-256 digest"
case "$expected" in
  *[!0-9A-Fa-f]*) fail "malformed SHA-256 digest" ;;
esac
if command -v sha256sum >/dev/null 2>&1; then
  actual=$(sha256sum "$tmp/$archive" | awk '{print $1}')
else
  actual=$(shasum -a 256 "$tmp/$archive" | awk '{print $1}')
fi
[ "$expected" = "$actual" ] || fail "checksum mismatch"

attestation_mode=${WATCHMAN_VERIFY_ATTESTATION:-${GPU_WATCHMAN_VERIFY_ATTESTATION:-auto}}
case "$attestation_mode" in
  auto|required)
    if command -v gh >/dev/null 2>&1 && gh attestation verify --help >/dev/null 2>&1; then
      printf '%s\n' "Verifying GitHub artifact attestation"
      GH_FORCE_TTY=0 gh attestation verify "$tmp/$archive" \
        --repo bas3line/watchman >/dev/null || fail "artifact attestation verification failed"
    elif [ "$attestation_mode" = required ]; then
      fail "a GitHub CLI with attestation support is required"
    else
      printf '%s\n' "Compatible GitHub CLI not found; continuing with SHA-256 verification"
    fi
    ;;
  disabled) ;;
  *) fail "WATCHMAN_VERIFY_ATTESTATION must be auto, required, or disabled" ;;
esac

member_names=$(tar -tzf "$tmp/$archive") || fail "cannot inspect release archive"
expected_members='gpu-watchman
README.md
CHANGELOG.md
SECURITY.md
LICENSE'
[ "$member_names" = "$expected_members" ] || fail "release archive has an unexpected member set"
member_listing=$(tar -tvzf "$tmp/$archive" gpu-watchman) || fail "cannot inspect release archive"
printf '%s\n' "$member_listing" | awk '
  NR == 1 && substr($1, 1, 1) == "-" { regular = 1; next }
  { regular = 0 }
  END { if (NR != 1 || !regular) exit 1 }
' >/dev/null || fail "archive must contain exactly one regular gpu-watchman file"
tar -xzf "$tmp/$archive" -C "$tmp" gpu-watchman
[ -f "$tmp/gpu-watchman" ] || fail "archive does not contain gpu-watchman"
[ ! -L "$tmp/gpu-watchman" ] || fail "archive contains a symbolic-link gpu-watchman"

mkdir -p "$INSTALL_DIR"
install -m 0755 "$tmp/gpu-watchman" "$INSTALL_DIR/watchman"
printf '%s\n' "Installed $INSTALL_DIR/watchman"
printf '%s\n' "Run: watchman version"
