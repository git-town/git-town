#!/usr/bin/env bash

# Performs initial configuration before running any Git Town command,
# unless the `--bypass-automatic-configuration` option is passed (used by git-town)


# Migrate old configuration (Git Town v0.2.2 and lower)
if [[ -f ".main_branch_name" ]]; then
  store_configuration main-branch-name "$(cat .main_branch_name)"
  rm .main_branch_name
fi


export main_branch_name=$(get_configuration main-branch-name)
export non_feature_branch_names=$(get_configuration non-feature-branch-names)


# Bypass the configuration if requested by caller (e.g. git-town)
if [[ $@ =~ --bypass-automatic-configuration ]]; then
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
    echo_and_exit_with_error
  fi
  echo
  store_configuration main-branch-name "$main_branch_name"
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
  store_configuration non-feature-branch-names "$non_feature_branch_names"
fi
