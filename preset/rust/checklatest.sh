#!/bin/bash
set -euo pipefail
curl https://static.rust-lang.org/dist/channel-rust-stable.toml --no-progress-meter | dasel -r toml -w json | jq -r '.pkg.rust.version|split(" ")[0]'