#!/usr/bin/env bash


function create_pull_request {
  local repository=$1
  local branch=$2
  local parent_branch=$3

  local to_compare
  if [ "$parent_branch" = "$MAIN_BRANCH_NAME" ]; then
    # Allow GitLab to redirect to the proper place if this repository is a fork
    to_compare="$branch"
  else
    to_compare="$parent_branch...$branch"
  fi

  open_browser "https://gitlab.com/$repository/compare/$to_compare?expand=1"
}


function show_repo {
  local repository=$1

  open_browser "https://gitlab.com/$repository"
}
