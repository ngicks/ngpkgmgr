#!/bin/bash
set -euo pipefail
curl -sSI -D - https://github.com/neovim/neovim/releases/latest -o /dev/null | grep -e '^location:' | sed 's/location: https:\/\/github.com\/neovim\/neovim\/releases\/tag\/v//'