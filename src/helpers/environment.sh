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


function git_town_revision {
  local base_dir=$(dirname "${BASH_SOURCE[0]}")
  local git_hash=$(git -C "$base_dir" rev-parse --short HEAD)
  local git_date=$(git -C "$base_dir" show -s --format=%cD HEAD | cut -d ' ' -f 2-4)

  echo "(${git_date}, ${git_hash})"
}


function is_git_town_installed_manually {
  local base_dir=$(dirname "${BASH_SOURCE[0]}")

  if [ "$(git -C "$base_dir" remote -v | grep -c "Originate/git-town")" == 0 ]; then
    echo false
  else
    echo true
  fi
}


if [[ ! "$@" =~ --bypass-environment-checks ]]; then
  ensure_git_repository
fi
