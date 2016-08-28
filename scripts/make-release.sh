#! /bin/bash

set -e

gox -verbose
mkdir -p release
mv dotfiles-command_* release/
cd release
for bin in `ls`; do
    mv "$bin" dotfiles-command
    zip "${bin}.zip" dotfiles-command
    rm dotfiles-command
done
cd -
