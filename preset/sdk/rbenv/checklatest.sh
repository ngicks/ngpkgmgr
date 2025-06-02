#!/bin/bash
curl -sSI -D -  https://github.com/rbenv/rbenv/releases/latest -o /dev/null | grep -e '^location:' | sed 's/location: https:\/\/github.com\/rbenv\/rbenv\/releases\/tag\/v//'
