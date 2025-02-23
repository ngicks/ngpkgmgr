#!/bin/bash
set -e
set -o pipefail
pushd /tmp
curl -LO https://go.dev/dl/go${VER}.${OS}-${ARCH}.tar.gz
sudo tar -C /usr/local -xzf go${VER}.${OS}-${ARCH}.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin/:$(/usr/local/go/bin/go env GOPATH)/bin' >> ~/home.sh
popd