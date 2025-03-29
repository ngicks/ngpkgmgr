#!/bin/bash
set -euo pipefail
pushd /tmp
curl -LO https://go.dev/dl/go${VER}.${OS}-${ARCH}.tar.gz
rm -rf ~/.local/go
tar -C ~/.local -xzf go${VER}.${OS}-${ARCH}.tar.gz
popd