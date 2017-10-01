#!/usr/bin/env bash

set -euET -o pipefail

main() {
  if [[ -t 1 ]]; then
    read -p "This will delete all uncommitted state. Are you sure? " answer
    case $answer in
      y|Y|yes) ;;
      n|N|no) return 0 ;;
      *) echo "Unparsable response $answer" >&2; exit 1 ;;
    esac
  fi
  git reset --hard
  git clean -f -d
  git submodule update --force --checkout --recursive
}

main
