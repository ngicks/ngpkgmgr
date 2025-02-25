#!/bin/bash
set -euo pipefail
pushd /tmp
curl -LO https://go.dev/dl/go${VER}.${OS}-${ARCH}.tar.gz
tar -C ~/.local/ -xzf go${VER}.${OS}-${ARCH}.tar.gz
popd
# in case it is not created.
mkdir ~/go -p
cp $(dirname $0)/env.sh ~/go/
echo '. ~/go/env.sh' >> ~/home.sh