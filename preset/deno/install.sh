#!/bin/bash
set -euo pipefail
curl -fsSL https://deno.land/install.sh | sh
echo 'export PATH=~/.deno/bin:$PATH' >> ~/home.sh