#!/bin/bash
set -euo pipefail

oss=(
  "linux"
  "darwin"
  "windows"
)

archs=(
  "amd64"
  "arm64"
)

for os in "${oss[@]}" ; do
  for arch in "${archs[@]}" ; do
    mkdir -p ./prebuilt/${os}-${arch}/
    go build -o ./prebuilt/${os}-${arch}/ .
    for p in ./cmd/*; do
      go build -trimpath -o ./prebuilt/${os}-${arch}/ $p
    done
  done
done
