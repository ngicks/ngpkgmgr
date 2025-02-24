#!/bin/bash
set -euo pipefail
curl -sSI -D - https://github.com/astral-sh/uv/releases/latest -o /dev/null | grep -e '^location:' | sed 's/location: https:\/\/github.com\/astral-sh\/uv\/releases\/tag\///'