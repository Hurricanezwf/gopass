#!/bin/bash

VERSION=0.0.4
SAVEDIR=release
OS=(linux darwin)
ARCH=(amd64 386)

set -e

rm -rf ./release

for os in ${OS[@]}
do
    for arch in ${ARCH[@]}
    do
        echo "start build $os $arch..."
        dir="${SAVEDIR}/gopass_v${VERSION}_${os}_${arch}"
        rm -rf ./$dir 
        mkdir -p ./$dir
        GOOS=$os GOARCH=$arch go build -o ./$dir/gopass
        cp ./conf.ini ./$dir/
        tar -zcvf ${dir}.tar.gz ./${dir}
        rm -rf ./${dir}
    done
done
