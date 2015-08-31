#!/usr/bin/env bash


# Exit if the current branch does not have extractable commits
function ensure_has_extractable_commits {
  local current_branch_name=$(get_current_branch_name)
  if [ "$(has_extractable_commits "$current_branch_name")" == false ]; then
    echo_error_header
    echo_error "The branch '$current_branch_name' has no extractable commits."
    exit_with_error newline
  fi
}


# Determines whether the given branch has extractable commits
function has_extractable_commits {
  local branch_name=$1
  if [ "$(git log --oneline "$MAIN_BRANCH_NAME..$branch_name" | wc -l | tr -d ' ')" == 0 ]; then
    echo false
  else
    echo true
  fi
}


function revert_commit {
  local commit=$1
  run_command "git revert $commit"
}
