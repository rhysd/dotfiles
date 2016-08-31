#! /bin/bash

set -e

gox -verbose
mkdir -p release
mv dotfiles_* release/
cd release
for bin in `ls`; do
    mv "$bin" dotfiles
    zip "${bin}.zip" dotfiles
    rm dotfiles
done
cd -
