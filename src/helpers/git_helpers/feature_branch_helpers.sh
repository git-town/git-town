#!/usr/bin/env bash


# Creates and checkouts a new branch off the main branch with the given name
function create_and_checkout_feature_branch {
  create_and_checkout_branch "$1" "$MAIN_BRANCH_NAME"
}


# Exits if the supplied branch is not a feature branch
function ensure_is_feature_branch {
  local branch_name=$1
  local error_message=$2
  if [ "$(is_feature_branch "$branch_name")" == false ]; then
    echo_error_header
    echo_error "The branch '$branch_name' is not a feature branch. $error_message"
    exit_with_error newline
  fi
}


# Exits if the current branch is not a feature branch
function ensure_on_feature_branch {
  local error_message="$*"
  local branch_name=$(get_current_branch_name)
  ensure_is_feature_branch "$branch_name" "$error_message"
}


# Returns true if the current branch is a feature branch
function is_feature_branch {
  local branch_name=$1
  if [ "$branch_name" == "$MAIN_BRANCH_NAME" -o "$(is_non_feature_branch "$branch_name")" == true ]; then
    echo false
  else
    echo true
  fi
}


# Returns true if the current branch is the main branch
function is_main_branch {
  local branch_name=$1

  if [ "$branch_name" == "$MAIN_BRANCH_NAME" ]; then
    echo true
  else
    echo false
  fi
}


# Returns true if the given branch is a non-feature branch
function is_non_feature_branch {
  local branch_name=$1

  if echo "$NON_FEATURE_BRANCH_NAMES" | tr ',' '\n' | grep -q "^ *$branch_name *$"; then
    echo true
  else
    echo false
  fi
}
