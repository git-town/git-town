#!/bin/bash

# Helpers methods that require a specific status or error


# Exits if the repository has a branch with the given name
function ensure_does_not_have_branch {
  local branch_name=$1
  if [ "$(has_branch "$branch_name")" = true ]; then
    echo_error_header
    echo_error "A branch named '$branch_name' already exists"
    exit_with_error
  fi
}



# Exits if the repository does not have a branch with the given name
function ensure_has_branch {
  local branch_name=$1
  if [ "$(has_branch "$branch_name")" == false ]; then
    echo_error_header
    echo_error "There is no branch named '$branch_name'."
    exit_with_error
  fi
}


# Exits if the supplied branch is not a feature branch
function ensure_is_feature_branch {
  local branch_name=$1
  local error_message=$2
  if [ "$(is_feature_branch "$branch_name")" == false ]; then
    error_is_not_feature_branch

    echo_error_header
    echo_error "The branch '$branch_name' is not a feature branch. $error_message"
    exit_with_error
  fi
}


# Exits if there are unresolved conflicts
function ensure_no_conflicts {
  if [ "$(has_conflicts)" == true ]; then
    echo_error_header
    echo_error "$*"
    exit_with_error
  fi
}


# Exists if there are uncommitted changes
function ensure_no_open_changes {
  if [ "$(has_open_changes)" == true ]; then
    error_has_open_changes

    echo_error_header
    echo_error "$*"
    exit_with_error
  fi
}


# Exits if the current branch is not a feature branch
function ensure_on_feature_branch {
  local error_message="$*"
  local branch_name=$(get_current_branch_name)
  ensure_is_feature_branch "$branch_name" "$error_message"
}
