#!/bin/bash
set -euo pipefail
curl https://go.dev/VERSION?m=text --no-progress-meter | head -n 1 | sed 's/go//'