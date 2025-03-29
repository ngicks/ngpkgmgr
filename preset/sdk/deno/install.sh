#!/bin/bash
set -euo pipefail
curl -fsSL https://deno.land/install.sh | sh
echo 'export PATH=$HOME/.deno/bin:$PATH' >> ${PROFILE_SH}
