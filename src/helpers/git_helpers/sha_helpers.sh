#!/usr/bin/env bash


# Returns the sha of the current branch
function current_sha {
  sha_of_branch HEAD
}


# Returns the SHA that the given branch points to
function sha_of_branch {
  local branch_name=$1
  git rev-parse "$branch_name"
}


# Resets the current branch to the commit described by the given SHA
function reset_to_sha {
  local sha=$1
  local hard=$2
  if [ "$sha" != "$(current_sha)" ]; then
    if [ -n "$hard" ]; then
      run_git_command "git reset --hard $sha"
    else
      run_git_command "git reset $sha"
    fi
  fi
}
