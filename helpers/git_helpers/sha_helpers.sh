#!/bin/bash


# Returns the SHA that the given branch points to
function sha_of_branch {
  local branch_name=$1
  git rev-parse "$branch_name"
}


# Resets the current branch to the commit described by the given SHA
function reset_to_sha {
  local sha=$1
  run_command "git reset $sha"
}
