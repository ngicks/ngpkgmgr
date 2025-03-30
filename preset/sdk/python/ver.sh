#!/bin/bash
set -euo pipefail
~/.local/uv_global/bin/python3 --version | sed 's/Python //'
