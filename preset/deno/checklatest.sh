#!/bin/bash
set -e
set -o pipefail
curl -sSI -D - https://github.com/denoland/deno/releases/latest -o /dev/null | grep -e '^location:' | sed 's/location: https:\/\/github.com\/denoland\/deno\/releases\/tag\/v//'