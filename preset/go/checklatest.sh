#!/bin/bash
set -e
set -o pipefail
curl https://go.dev/VERSION?m=text --no-progress-meter | head -n 1 | sed 's/go//'