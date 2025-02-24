#!/bin/bash
set -euo pipefail
go version | awk '{print $3}' | sed 's/go//'