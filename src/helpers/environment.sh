#!/usr/bin/env bash

function is_git_repository {
  if git rev-parse > /dev/null 2>&1; then
    echo true
  else
    echo false
  fi
}


function ensure_git_repository {
  if [ "$(is_git_repository)" == false ]; then
    echo_inline_error "This is not a git repository."
    exit_with_error
  fi
}


function gittown_dev_version {
  local base_dir=$(dirname "${BASH_SOURCE[0]}")


  if [ "$(git -C "$base_dir" remote -v | grep -c "https://github.com/Homebrew/homebrew.git (fetch)")" == 0 ]; then
    local git_hash=$(git -C "$base_dir" rev-parse --short HEAD)
    local git_date=$(git -C "$base_dir" show -s --format=%cD HEAD | cut -d ' ' -f 2-4)

    echo " (${git_date}, ${git_hash})"
  fi
}


if [[ ! "$@" =~ --bypass-environment-checks ]]; then
  ensure_git_repository
fi
