#!/bin/sh
set -eu
umask 077

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
VERSION=v0.8.2
MOCK_BIN=$ROOT/scripts/testdata/watchman-bin
RELEASE_ROOT=$ROOT/public/watchman/releases/$VERSION

fail() {
  printf '%s\n' "watchman installer test: $*" >&2
  exit 1
}

[ -d "$RELEASE_ROOT" ] || fail "missing $RELEASE_ROOT"
grep -F "VERSION=\${WATCHMAN_VERSION:-\${GPU_WATCHMAN_VERSION:-$VERSION}}" "$ROOT/install.sh" >/dev/null ||
  fail "installer default is not $VERSION"

test_root=$(mktemp -d "${TMPDIR:-/tmp}/watchman-installer-test.XXXXXX")
trap 'rm -rf "$test_root"' EXIT INT TERM

host_os=$(uname -s | tr '[:upper:]' '[:lower:]')
host_arch=$(uname -m)
case "$host_arch" in
  x86_64|amd64) host_arch=amd64 ;;
  aarch64|arm64) host_arch=arm64 ;;
  *) host_arch=unsupported ;;
esac
host_suffix=${host_os}_${host_arch}

for target in \
  'Darwin x86_64 darwin_amd64' \
  'Darwin arm64 darwin_arm64' \
  'Linux x86_64 linux_amd64' \
  'Linux aarch64 linux_arm64'
do
  set -- $target
  uname_s=$1
  uname_m=$2
  suffix=$3
  install_dir=$test_root/$suffix/bin
  expected_dir=$test_root/$suffix/expected
  mkdir -p "$install_dir" "$expected_dir"

  archive=$RELEASE_ROOT/gpu-watchman_${suffix}.tar.gz
  checksum=$archive.sha256
  [ -f "$archive" ] || fail "missing $(basename "$archive")"
  [ -f "$checksum" ] || fail "missing $(basename "$checksum")"

  members=$(tar -tzf "$archive") || fail "cannot inspect $(basename "$archive")"
  expected_members='gpu-watchman
README.md
CHANGELOG.md
SECURITY.md
LICENSE'
  [ "$members" = "$expected_members" ] || fail "$(basename "$archive") has an unexpected member set"
  tar -xzf "$archive" -C "$expected_dir" gpu-watchman

  PATH=$MOCK_BIN:$PATH \
    MOCK_UNAME_S=$uname_s \
    MOCK_UNAME_M=$uname_m \
    MOCK_WATCHMAN_REGISTRY_ROOT=$ROOT/public \
    WATCHMAN_BASE_URL=https://registry.test \
    WATCHMAN_VERIFY_ATTESTATION=disabled \
    WATCHMAN_INSTALL_DIR=$install_dir \
    sh "$ROOT/install.sh" >/dev/null

  [ -x "$install_dir/watchman" ] || fail "installed $suffix binary is not executable"
  cmp "$expected_dir/gpu-watchman" "$install_dir/watchman" >/dev/null ||
    fail "installed $suffix binary differs from the verified archive member"
  if [ "$suffix" = "$host_suffix" ]; then
    installed_version=$($install_dir/watchman version) || fail "native $suffix binary did not run"
    [ "$installed_version" = "${VERSION#v}" ] ||
      fail "native $suffix binary reports version $installed_version, expected ${VERSION#v}"
  fi
done

if PATH=$MOCK_BIN:$PATH \
  MOCK_UNAME_S=Linux \
  MOCK_UNAME_M=x86_64 \
  MOCK_WATCHMAN_REGISTRY_ROOT=$ROOT/public \
  MOCK_WATCHMAN_CORRUPT_ARCHIVE=1 \
  WATCHMAN_BASE_URL=https://registry.test \
  WATCHMAN_VERIFY_ATTESTATION=disabled \
  WATCHMAN_INSTALL_DIR=$test_root/corrupt/bin \
  sh "$ROOT/install.sh" >"$test_root/corrupt.out" 2>&1
then
  fail "installer accepted an archive with a checksum mismatch"
fi
grep -F 'checksum mismatch' "$test_root/corrupt.out" >/dev/null ||
  fail "checksum mismatch did not return the expected diagnostic"

if WATCHMAN_BASE_URL=http://registry.test \
  WATCHMAN_INSTALL_DIR=$test_root/insecure/bin \
  sh "$ROOT/install.sh" >"$test_root/insecure.out" 2>&1
then
  fail "installer accepted an insecure registry URL"
fi
grep -F 'WATCHMAN_BASE_URL must use HTTPS' "$test_root/insecure.out" >/dev/null ||
  fail "insecure URL did not return the expected diagnostic"

printf '%s\n' "Watchman installer checks passed for $VERSION on macOS/Linux amd64/arm64"
