#!/usr/bin/env bash


# Checks out the branch with the given name (if not alread checked out)
function checkout {
  local branch_name=$1
  if [ "$(get_current_branch_name)" != "$branch_name" ]; then
    run_command "git checkout $branch_name"
  fi
}


# Checkout a branch without this showing up in the output
function checkout_silently {
  local branch_name=$1
  run_command_silently "git checkout $branch_name"
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
