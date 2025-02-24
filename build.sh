#!/bin/bash
set -euo pipefail
mkdir -p ./prebuilt/$(go env GOOS)-$(go env GOARCH)/
go build -o ./prebuilt/$(go env GOOS)-$(go env GOARCH)/ .
for p in ./cmd/*; do
    go build -o ./prebuilt/$(go env GOOS)-$(go env GOARCH)/ $p
done