#!/bin/bash
set -euo pipefail
curl https://get.volta.sh | bash -s -- --skip-setup
echo export 'VOLTA_HOME="$HOME/.volta"' >> ~/home.sh
echo export 'PATH="$VOLTA_HOME/bin:$PATH"' >> ~/home.sh