#! /bin/bash

set -e

git config --global url.https://github.com/.insteadOf git@github.com:

if [[ "$TRAVIS_OS_NAME" == "osx" ]]; then
    brew update
    brew upgrade go
    go get -t -d -v ./...
    go test -v ./src/
    go test ./
else
    go get github.com/axw/gocov/gocov
    go get github.com/mattn/goveralls
    if ! go get code.google.com/p/go.tools/cmd/cover; then go get golang.org/x/tools/cmd/cover; fi
    go get -t -d -v ./...
    go vet
    cd src/ && go vet && cd -
    go test ./
    go test -v -coverprofile=coverage.out ./src/
    $HOME/gopath/bin/goveralls -coverprofile coverage.out -service=travis-ci -repotoken $COVERALLS_TOKEN
fi

