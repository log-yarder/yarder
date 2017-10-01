#!/usr/bin/env bash

set -euET -o pipefail

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
  local err=0
  ("$@") || err=$?
  travis_fold end "$name"
  return $err
}

run_goimports() {
  go get -v golang.org/x/tools/cmd/goimports || return $?
  local err=0
  for d in $(go list -f '{{.Dir}}' ./... | grep -v '/vendor/'); do
    echo "goimports $d/*.go"
    test -z "$(goimports -d "$d"/*.go | tee /dev/stderr)" || err=1
  done
  return $err
}

run_gotest() {
  local err=0
  for p in $(go list ./... | grep -v '/vendor/'); do
    echo "go test -v $p"
    go test -v "$p" || err=$?
  done
  return $err
}

run_govet() {
  local err=0
  for p in $(go list ./... | grep -v '/vendor/'); do
    echo "go vet $p"
    go vet "$p" || err=$?
  done
  return $err
}

main() {
  local err=0
  run "goimports" run_goimports || err=$?
  run "go-test" run_gotest || err=$?
  run "go-vet" run_govet || err=$?
  return $err
}

main
