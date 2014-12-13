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

  # update $main_branch_name and $non_feature_branch_names accordingly
  if [ $? == '0' ]; then
    if [ "$config_setting_name" == "main-branch-name" ]; then
      main_branch_name="$value"
    elif [ "$config_setting_name" == "non-feature-branch-names" ]; then
      non_feature_branch_names="$value"
    fi
  fi
}


# Persists the main branch configuration
function store_main_branch_name_with_confirmation_text {
  store_configuration main-branch-name "$1"
  echo "main branch stored as '$1'"
}


# Persists the non-feature branch configuration
function store_non_feature_branch_names_with_confirmation_text {
  store_configuration non-feature-branch-names "$1"
  echo "non-feature branches stored as '$1'"
}


# Update old configuration to new one if it exists
if [[ -f ".main_branch_name" ]]; then
  store_configuration main-branch-name "$(cat .main_branch_name)"
  rm .main_branch_name
fi


main_branch_name=$(get_configuration main-branch-name)
non_feature_branch_names=$(get_configuration non-feature-branch-names)


# Bypass the configuration if requested by caller (e.g. git-town)
if [[ $1 == "--bypass-automatic-configuration" ]]; then
  return 0
fi


# Ask and store main-branch-name, if it isn't known yet.
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
  echo
  store_main_branch_name_with_confirmation_text "$main_branch_name"
fi

# Ask and store non-feature-branch-names, if needed
if [[ $? == '1' ]]; then
  echo
  echo "Git Town supports non-feature branches like 'release' or 'production'."
  echo "These branches cannot be shipped and do not merge $main_branch_name when syncing."
  echo "Please enter the names of all your non-feature branches as a comma separated list."
  echo "Example: 'qa, production'"
  read non_feature_branch_names
  echo
  store_non_feature_branch_names_with_confirmation_text "$non_feature_branch_names"
fi
