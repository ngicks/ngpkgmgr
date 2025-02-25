#!/bin/bash
set -euo pipefail
~/.local/bin/uv python install ${VER}
# uv creates ~/.local/bin/ dir
ln -sf $(~/.local/bin/uv python find ${VER}) ~/.local/bin/python3
ln -sf python3 ~/.local/bin/python