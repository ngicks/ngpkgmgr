#!/bin/bash

rbenv="$(command -v rbenv ~/.local/rbenv/bin/rbenv | head -1)"

if [ -n "$rbenv" ]; then
  echo "rbenv already seems installed in \`$rbenv'."
  cd "${rbenv%/*}"

  if git remote -v 2>/dev/null | grep -q rbenv; then
    echo "Trying to update with git..."
    git pull --tags origin master
    cd ..
  fi
else
  echo "Installing rbenv with git..."
  mkdir -p ~/.local/rbenv
  cd ~/.local/rbenv
  git init
  git remote add -f -t master origin https://github.com/rbenv/rbenv.git
  git checkout -b master origin/master
  git reset --hard tags/v${VER}
  rbenv=~/.local/rbenv/bin/rbenv
fi

ruby_build="$(command -v "$rbenv_root"/plugins/*/bin/rbenv-install rbenv-install | head -1)"

if [ -n "$ruby_build" ]; then
  echo "\`rbenv install' command already available in \`$ruby_build'."
  cd "${ruby_build%/*}"
  if git remote -v 2>/dev/null | grep -q ruby-build; then
    echo "Trying to update with git..."
    git pull origin master
  fi
else
  echo "Installing ruby-build with git..."
  mkdir -p ~/.local/rbenv/plugins
  git clone https://github.com/rbenv/ruby-build.git ~/.local/rbenv/plugins/ruby-build
fi

mkdir -p ~/.local/rbenv/cache

echo 'export RBENV_ROOT=$HOME/.local/rbenv' >> ${PROFILE_SH}
echo 'eval "$($HOME/.local/rbenv/bin/rbenv init - --no-rehash bash)"' >> ${PROFILE_SH}
