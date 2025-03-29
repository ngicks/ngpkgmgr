#!/bin/bash
~/.local/bin/uv python install
uv venv ~/.local/uv_global
echo 'export PATH="$HOME/.local/uv_global/bin:$PATH"' >> ~/home.sh

