#!/bin/bash
~/.local/rbenv/bin/rbenv --version | sed 's/rbenv //' | sed 's/-[a-z0-9\-]\+//'
