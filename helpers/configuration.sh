#!/bin/bash

# Helper methods for dealing with configuration.

# Returns git-town configuration
function get_configuration {
  local config=$1
  git config "git-town.$config"
}


# Persists git-town configuration
function store_configuration {
  local config=$1
  local value=$2
  git config "git-town.$config" "$value"
}


# Persists the main branch configuration
function store_main_branch_name {
  store_configuration "main-branch-name" "$1"
}


# Persists the non-feature branch configuration
function store_non_feature_branch_names {
  store_configuration "non-feature-branch-names" "$1"
}


# Update old configuration to new one if it exists
if [[ -f ".main_branch_name" ]]; then
  store_main_branch_name "$(cat .main_branch_name)"
  rm .main_branch_name
fi


# Read main branch name from config, ask and store it if it isn't known yet.
main_branch_name=$(git config git-town.main-branch-name)
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
non_feature_branch_names=$(git config git-town.non-feature-branch-names)
if [[ $? == '1' ]]; then
  echo "Git Town supports non-feature branches like 'release' or 'production'."
  echo "These branches cannot be shipped and do not merge $main_branch_name when syncing."
  echo "Please enter the names of all your non-feature branches as a comma seperated list."
  echo "Example: 'qa, production'"
  read non_feature_branch_names
  store_non_feature_branch_names "$non_feature_branch_names"
  echo "non-feature branches stored as '$non_feature_branch_names'"
fi
