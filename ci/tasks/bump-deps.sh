#!/bin/bash

set -e

git clone bosh-cpi-src-in bosh-cpi-src-out

export GOPATH=$PWD/bosh-cpi-src-out

cd bosh-cpi-src-out/src/bosh-google-cpi

#intentionally cause an explicit commit if the underlying go version in our compiled Dockerfile changes.
#assume the go-dep-bumper and bosh-utils are bumping at the same cadence.
NEW_GO_MINOR=$(go version | sed 's/go version go1\.\([0-9]\+\)\..*$/\1/g')
CURRENT_GO_MINOR=$(cat go.mod | grep -E "^go 1.*" | sed "s/^\(go 1.\)\([0-9]\+\)/\2/")
sed -i "s/^\(go 1.\)\([0-9]\+\)/\1$NEW_GO_MINOR/" go.mod

go get -u ./...
go mod tidy
go mod vendor

if [ "$(git status --porcelain)" != "" ]; then
  git status
  git add vendor go.sum go.mod
  git config user.name "CI Bot"
  git config user.email "cf-bosh-eng@pivotal.io"
  if [ $CURRENT_GO_MINOR == $NEW_GO_MINOR ]; then
    git commit -m "Update vendored dependencies"
  else
    git commit -m "Bump to go version 1.$NEW_GO_MINOR\n\n- (and update vendored dependencies)"
  fi
fi
