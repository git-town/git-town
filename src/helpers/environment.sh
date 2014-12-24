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
    echo "This is not a git repository."
    exit_with_error
  fi
}


# Bypass the environment checks
if [[ $@ =~ "--bypass-environment-checks" ]]; then
  return 0
fi

ensure_git_repository

export initial_branch_name=$(get_current_branch_name)
export initial_open_changes=$(has_open_changes)
