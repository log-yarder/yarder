#!/usr/bin/env bash

set -euET -o pipefail

main() {
  git submodule update --force --checkout --recursive
  go get -v golang.org/x/tools/cmd/goimports
  for d in $(go list -f '{{.Dir}}' ./... | grep -v '/vendor/'); do
    goimports -w "$d"/*.go
  done
}

main
