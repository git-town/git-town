#!/bin/bash

# Helper methods that return some output


# Returns the current branch name
function get_current_branch_name {
  git rev-parse --abbrev-ref HEAD
}


# Returns the names of local branches that have been merged into main
function local_merged_branches {
  git branch --merged "$main_branch_name" | tr -d ' ' | sed 's/\*//g'
}


# Returns the names of remote branches that have been merged into main
function remote_merged_branches {
  git branch -r --merged "$main_branch_name" | grep -v HEAD | tr -d ' ' | sed 's/origin\///g'
}


# Returns the SHA that the given branch points to
function sha_of_branch {
  local branch_name=$1
  git rev-parse "$branch_name"
}
