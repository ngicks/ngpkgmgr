#!/bin/bash
~/.cargo/bin/rustc +stable --version | sed 's/rustc\s//' | sed 's/ (.*)//'