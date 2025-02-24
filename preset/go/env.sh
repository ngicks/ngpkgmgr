#!/bin/bash

gobin=/usr/local/go/bin/go

export PATH=$($gobin env GOROOT)/bin:$PATH

if [[ -n $($gobin env GOBIN) ]]; then
    export PATH=$($gobin env GOBIN):$PATH
else
    export PATH=$($gobin env GOPATH)/bin:$PATH
fi
