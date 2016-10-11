#! /bin/bash

set -e

git config --global url.https://github.com/.insteadOf git@github.com:
echo -e "Host github.com\n\tVerifyHostKeyDNS no\n" >> ~/.ssh/config

if [[ "$TRAVIS_OS_NAME" == "osx" ]]; then
    brew update
    set +e
    brew upgrade go
    set -e
    go get -t -d -v ./...
    set +e
    go test ./
    set -e
    go test -v ./src/
else
    go get github.com/axw/gocov/gocov
    go get github.com/mattn/goveralls
    go get golang.org/x/tools/cmd/cover
    go get -t -d -v ./...
    go vet
    set +e
    go test ./
    set -e
    cd src/ && go vet && cd -
    go test -v -coverprofile=coverage.out ./src/
    $HOME/gopath/bin/goveralls -coverprofile coverage.out -service=travis-ci -repotoken $COVERALLS_TOKEN
fi

