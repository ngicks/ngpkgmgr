#!/bin/bash
curl -LsSf https://astral.sh/uv/install.sh | NO_MODIFY_PATH=1 sh
echo '. "$HOME/.local/bin/env"' >> ${PROFILE_SH}
