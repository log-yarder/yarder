#!/usr/bin/env bash

set -euET -o pipefail

die() {
  echo "$*"
  exit 1
}

travis_fold() {
  local type="$1"
  local group="$2"
  if [[ "${TRAVIS:-}" != true ]]; then
    return
  fi
  echo "travis_fold:$type:$group"
}

run() {
  local name="$1"
  shift
  travis_fold start "$name"
  echo "Running $*"
  err=0
  ("$@") || err=$?
  travis_fold end "$name"
  return $err
}

run_goimports() {
  err=0
  for d in $(go list -f '{{.Dir}}' ./...); do
    echo "goimports $d/*.go"
    test -z "$(goimports -d "$d"/*.go | tee /dev/stderr)" || err=1
  done
  return $err
}

run "get-goimports" go get -v golang.org/x/tools/cmd/goimports
run "goimports" run_goimports
run "go-test" go test -v ./...
run "go-vet" go vet ./...
