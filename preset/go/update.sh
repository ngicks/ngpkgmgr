#!/bin/bash
set -e
set -o pipefail
pushd /tmp
curl -LO https://go.dev/dl/go${VER}.${OS}-${ARCH}.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go${VER}.${OS}-${ARCH}.tar.gz
popd