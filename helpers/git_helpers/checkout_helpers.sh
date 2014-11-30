#!/bin/bash


# Checks out the branch with the given name.
#
# Skips this operation if the requested branch
# is already checked out.
function checkout_branch {
  local branch_name=$1
  if [ ! "$(get_current_branch_name)" = "$branch_name" ]; then
    run_command "git checkout $branch_name"
  fi
}


# Checks out the main development branch in Git.
#
# Skips the operation if we already are on that branch.
function checkout_main_branch {
  checkout_branch "$main_branch_name"
}


# Cuts a new branch off the given parent branch, and checks it out.
function create_and_checkout_branch {
  local new_branch_name=$1
  local parent_branch_name=$2
  run_command "git checkout -b $new_branch_name $parent_branch_name"
}


# Creates a new feature branch with the given name.
#
# The feature branch is cut off the main development branch.
function create_and_checkout_feature_branch {
  create_and_checkout_branch "$1" "$main_branch_name"
}
