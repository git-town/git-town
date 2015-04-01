#!/usr/bin/env bash

function is_git_repository {
  if git rev-parse > /dev/null 2>&1; then
    echo true
  else
    echo false
  fi
}


function is_gittown_nightly {
  local base_dir=$( dirname "${BASH_SOURCE[0]}" )

  if git -C "$base_dir" remote -v > /dev/null 2>&1 | grep "https://github.com/Homebrew/homebrew.git (fetch)"; then
    echo false
  else
    echo true
  fi
}


function ensure_git_repository {
  if [ "$(is_git_repository)" == false ]; then
    echo_inline_error "This is not a git repository."
    exit_with_error
  fi
}


function gittown_nightly_version {
  local base_dir=$( dirname "${BASH_SOURCE[0]}" )
  local git_hash=$(git -C "$base_dir" rev-parse --short HEAD)
  local git_date=$(git -C "$base_dir" --no-pager show -s --format=%ci HEAD | cut -d ' ' -f1 | tr -d '-')

  if [ "$(is_gittown_nightly)" == true ]; then
    echo ".${git_date}-nightly (${git_hash})"
  fi
}

if [[ ! "$@" =~ --bypass-environment-checks ]]; then
  ensure_git_repository
fi
