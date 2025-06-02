#!/bin/bash
cd ~/.local/rbenv
git fetch --all
git checkout master
git reset --hard tags/v${VER}
