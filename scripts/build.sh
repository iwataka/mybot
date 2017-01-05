#!/usr/bin/env bash

cd $(dirname $0)/..

name="mybot"

build() {
    os=$1
    arch=$2
    echo Building for ${os}-${arch}...
    set +e
    mkdir -p "bin/${os}-${arch}" 2> /dev/null
    set -e
    exe="${name}"
    if [[ ${os} = "windows" ]]; then
        exe="${name}".exe
    fi
    GOOS=${os} GOARCH=${arch} go build -o "bin/${os}-${arch}/${exe}"
}

build linux 386
build linux amd64
build darwin 386
build darwin amd64
build windows 386
build windows amd64
