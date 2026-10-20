#!/bin/sh
set -eu
version=${1:-dev}
mkdir -p dist
for target in linux/386 linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64; do
    os=${target%/*}
    arch=${target#*/}
    name="iproxy-${os}-${arch}"
    if [ "$os" = windows ]; then name="${name}.exe"; fi
    CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" GO386=softfloat go build \
        -trimpath -ldflags "-s -w -X main.version=$version" -o "dist/$name" .
done
(
    cd dist
    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum iproxy-* > SHA256SUMS
    else
        shasum -a 256 iproxy-* > SHA256SUMS
    fi
)
