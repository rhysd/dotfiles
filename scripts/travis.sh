#! /bin/bash

set -e

if [[ "$TRAVIS_OS_NAME" == "osx" ]]; then
    brew update
    brew upgrade go
    go get -t -d -v ./...
    go test -v
else
    go get github.com/axw/gocov/gocov
    go get github.com/mattn/goveralls
    if ! go get code.google.com/p/go.tools/cmd/cover; then go get golang.org/x/tools/cmd/cover; fi
    go get -t -d -v ./...
    go test -v -coverprofile=coverage.out
    $HOME/gopath/bin/goveralls -coverprofile coverage.out -service=travis-ci -repotoken $COVERALLS_TOKEN
fi

