#!/bin/bash

set -euo pipefail

export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

TAG=v0.1.0

go build -a -ldflags '-extldflags "-static"' .

mkdir -p release
cd release/

OUTPUT_DIR="$GOOS-$GOARCH"

rm -rf "$OUTPUT_DIR"
mkdir "$OUTPUT_DIR"

cp ../LICENSE ../README.md ../mkharp "$OUTPUT_DIR"

TAR_FILE="mkharp-$TAG-$GOOS-$GOARCH.tar.gz"

tar czf "$TAR_FILE" "$OUTPUT_DIR"
sha256sum "$TAR_FILE" > "$TAR_FILE.sha256sum"

rm -rf "$OUTPUT_DIR"
