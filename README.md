# ngpkgmgr

A simple meta package manager which just stores shell commands for install / update / remove(WIP) pkg.

## prerequisits

- 7z (unzipping)
- jq (json manipulation)
- dasel (transofrming toml to json)
- ./cmd/picklatest (picking latest LTS(even major version) of nodejs)
- ./cmd/gobuildinfo (simple wrapper of `buildinfo.ReadFile`. As of Go 1.24, VCS info is embedded as build info.)

```
sudo apt update && sudo apt install -y p7zip-full jq dasel
mkdir -p ~/bin
cp ./prebuilt/linux-amd64/* ~/bin/
```