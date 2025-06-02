#!/bin/bash
~/.local/rbenv/bin/rbenv install -l | grep -e '^[0-9]' | tail -1
