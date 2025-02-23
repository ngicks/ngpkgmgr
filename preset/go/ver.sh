#!/bin/bash
set -e
set -o pipefail
go version | awk '{print $3}' | sed 's/go//'