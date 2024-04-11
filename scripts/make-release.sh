#! /bin/bash

set -e

gox -verbose -osarch="linux/amd64 linux/arm darwin/amd64 darwin/arm64 netbsd/amd64 openbsd/amd64 windows/amd64"
mkdir -p release
mv dotfiles_* release/
cd release
for bin in `ls`; do
    if [[ "$bin" == *windows* ]]; then
        command="dotfiles.exe"
    else
        command="dotfiles"
    fi
    mv "$bin" "$command"
    zip "${bin}.zip" "$command"
    rm "$command"
done
cd -
