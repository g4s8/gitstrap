#!/bin/sh
set -e

OWNER="g4s8"
NAME="gitstrap"
RELEASES_URL="https://github.com/$OWNER/$NAME/releases"
test -z "$TMPDIR" && TMPDIR="$(mktemp -d)"
TAR_FILE="$TMPDIR/$NAME.tar.gz"

last_version() {
  curl -sL -o /dev/null -w %{url_effective} "$RELEASES_URL/latest" | 
    rev | 
    cut -f1 -d'/'| 
    rev
}

download() {
  test -z "$VERSION" && VERSION="$(last_version)"
  test -z "$VERSION" && {
    echo "Unable to get $NAME version." >&2
    exit 1
  }
  rm -f "$TAR_FILE"
  local url="$RELEASES_URL/download/$VERSION/${NAME}_$(uname -s)_$(uname -m).tar.gz"
  curl -sL -o "$TAR_FILE" $url
}

download
tar -xzf "$TAR_FILE" -C "$TMPDIR"
rm -f $TAR_FILE
mv -i $TMPDIR $PWD/$NAME

