#! /bin/bash

set -e

gox -verbose
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
