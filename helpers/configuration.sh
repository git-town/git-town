#!/bin/bash

# Helper methods for dealing with configuration.

# Returns git-town configuration
function get_configuration {
  local config_setting_name=$1
  git config "git-town.$config_setting_name"
}


# Persists git-town configuration
function store_configuration {
  local config_setting_name=$1
  local value=$2
  git config "git-town.$config_setting_name" "$value"
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


# Exits if the current branch is not a feature branch
function ensure_on_feature_branch {
  local error_message="$*"
  local branch_name=$(get_current_branch_name)
  ensure_is_feature_branch "$branch_name" "$error_message"
}


# Returns true if the current branch is a feature branch
function is_feature_branch {
  local branch_name=$1
  if [ "$branch_name" == "$main_branch_name" -o "$(echo "$non_feature_branch_names" | tr ',' '\n' | grep -c "$branch_name")" == 1 ]; then
    echo false
  else
    echo true
  fi
}


# Persists the main branch configuration
function store_main_branch_name {
  store_configuration main-branch-name "$1"
}


# Persists the non-feature branch configuration
function store_non_feature_branch_names {
  store_configuration non-feature-branch-names "$1"
}


# Update old configuration to new one if it exists
if [[ -f ".main_branch_name" ]]; then
  store_main_branch_name "$(cat .main_branch_name)"
  rm .main_branch_name
fi


# Read main branch name from config, ask and store it if it isn't known yet.
main_branch_name=$(get_configuration main-branch-name)
if [[ -z "$main_branch_name" ]]; then
  echo "Please enter the name of the main dev branch (typically 'master' or 'development'):"
  read main_branch_name
  if [[ -z "$main_branch_name" ]]; then
    echo_error_header
    echo_error "You have not provided the name for the main branch."
    echo_error "This information is necessary to run this script."
    echo_error "Please try again."
    exit_with_error
  fi
  store_main_branch_name "$main_branch_name"
  echo
  echo "main branch stored as '$main_branch_name'."
fi


# Read non feature branch names from config, ask and store if needed
non_feature_branch_names=$(get_configuration non-feature-branch-names)
if [[ $? == '1' ]]; then
  echo "Git Town supports non-feature branches like 'release' or 'production'."
  echo "These branches cannot be shipped and do not merge $main_branch_name when syncing."
  echo "Please enter the names of all your non-feature branches as a comma seperated list."
  echo "Example: 'qa, production'"
  read non_feature_branch_names
  store_non_feature_branch_names "$non_feature_branch_names"
  echo "non-feature branches stored as '$non_feature_branch_names'"
fi
