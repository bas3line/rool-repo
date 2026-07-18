#!/bin/sh
set -eu

output=
url=
while [ "$#" -gt 0 ]; do
  case "$1" in
    --output)
      shift
      output=$1
      ;;
    --)
      shift
      url=$1
      ;;
    https://*) url=$1 ;;
  esac
  shift
done

[ -n "$url" ] || { printf '%s\n' 'mock curl: missing URL' >&2; exit 1; }
path=${url%%\?*}
case "$path" in
  */sandbox/install.sh)
    [ -n "$output" ] || { printf '%s\n' 'mock curl: missing output' >&2; exit 1; }
    cp "$SANDBOX_TEST_INSTALLER" "$output"
    ;;
  */sandbox/latest)
    if [ -n "$output" ]; then
      printf '%s\n' "$SANDBOX_TEST_VERSION" > "$output"
    else
      printf '%s\n' "$SANDBOX_TEST_VERSION"
    fi
    ;;
  */sandbox/releases/*)
    [ -n "$output" ] || { printf '%s\n' 'mock curl: missing output' >&2; exit 1; }
    name=${path##*/}
    cp "$SANDBOX_TEST_FIXTURE/$name" "$output"
    ;;
  *)
    printf '%s\n' "mock curl: unexpected URL $url" >&2
    exit 1
    ;;
esac
