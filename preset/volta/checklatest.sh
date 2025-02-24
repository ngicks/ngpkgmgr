#!/bin/bash
set -euo pipefail
curl -sSI -D - https://github.com/volta-cli/volta/releases/latest -o /dev/null | grep -e '^location:' | sed 's/location: https:\/\/github.com\/volta-cli\/volta\/releases\/tag\/v//'