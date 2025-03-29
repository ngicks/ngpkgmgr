#!/bin/bash
set -euo pipefail
curl https://registry.npmjs.org/node/ --no-progress-meter | jq '.versions | keys' | picklatest --even