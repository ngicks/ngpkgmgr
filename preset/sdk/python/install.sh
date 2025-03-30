#!/bin/bash
~/.local/bin/uv python install
~/.local/bin/uv venv ~/.local/uv_global
echo 'export PATH="$HOME/.local/uv_global/bin:$PATH"' >> ${PROFILE_SH}

