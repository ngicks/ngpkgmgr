#!/bin/bash
set -euo pipefail
nvim --version | head -n 1 | sed 's/NVIM v//'