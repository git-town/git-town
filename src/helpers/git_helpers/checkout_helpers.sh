#!/usr/bin/env bash


# Checks out the branch with the given name (if not alread checked out)
function checkout {
  local branch_name=$1
  if [ "$(get_current_branch_name)" != "$branch_name" ]; then
    run_command "git checkout $branch_name"
  fi
}


# Checks out the main branch (if not alread checked out)
function checkout_main_branch {
  checkout "$MAIN_BRANCH_NAME"
}


function undo_steps_for_checkout {
  echo "checkout $(get_current_branch_name)"
}


function undo_steps_for_checkout_main_branch {
  undo_steps_for_checkout
}
