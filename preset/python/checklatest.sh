#!/bin/bash
set -euo pipefail
curl https://endoflife.date/api/python.json --no-progress-meter | jq -r '.[0]["latest"]'