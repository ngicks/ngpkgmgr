#!/bin/bash
set -euo pipefail
curl https://get.volta.sh | bash -s -- --skip-setup
echo export 'VOLTA_HOME="$HOME/.volta"' >> ${PROFILE_SH}
echo export 'PATH="$VOLTA_HOME/bin:$PATH"' >> ${PROFILE_SH}
