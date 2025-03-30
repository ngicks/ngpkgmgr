#!/bin/bash
set -euo pipefail
pushd /tmp
curl -LO https://github.com/neovim/neovim/releases/download/v${VER}/nvim-${OS}-$(uname -m).tar.gz
rm -rf nvim-linux-x86_64
tar -xzf nvim-linux-x86_64.tar.gz
rm -rf ~/.local/nvim
mv ./nvim-linux-x86_64 ~/.local/nvim
popd

