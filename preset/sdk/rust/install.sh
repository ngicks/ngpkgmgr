#!/bin/bash
set -euo pipefail
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- --no-modify-path
echo '. "$HOME/.cargo/env"' >> ${PROFILE_SH}
