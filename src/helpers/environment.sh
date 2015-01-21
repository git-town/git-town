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
    exit 1
  fi
}


if [[ ! $@ =~ --bypass-environment-checks ]]; then
  ensure_git_repository
fi
