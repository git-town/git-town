#!/bin/bash


# Checks out the branch with the given name (if not alread checked out)
function checkout {
  local branch_name=$1
  if [ "$(get_current_branch_name)" != "$branch_name" ]; then
    run_command "git checkout $branch_name"
  fi
}


# Checks out the main branch (if not alread checked out)
function checkout_main_branch {
  checkout "$main_branch_name"
}
