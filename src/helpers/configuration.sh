#!/usr/bin/env bash

# Helper methods for dealing with configuration.

function echo_non_feature_branch_usage {
  echo_inline_usage 'git town non-feature-branches (--add | --remove) <branchname>'
}


# Add a new non-feature branch if possible, and show confirmation
function add_non_feature_branch {
  local branch_name=$1

  if [ -z "$branch_name" ]; then
    echo_inline_error "missing branch name"
    echo_non_feature_branch_usage
    exit 1
  elif [ "$(has_branch "$branch_name")" = false ]; then
    echo_inline_error "no branch named '$branch_name'"
    exit 1
  elif [ "$(is_non_feature_branch "$branch_name")" = true ]; then
    echo_inline_error "'$branch_name' is already a non-feature branch"
    exit 1
  else
    local new_branches=$(insert_string "$non_feature_branch_names" ',' "$branch_name")
    store_configuration non-feature-branch-names "$new_branches"
  fi
}

# Add or remove non-feature branch if possible, and show confirmation
function add_or_remove_non_feature_branches {
  local option=$1
  local branch_name=$2

  if [ "$option" == "--add" ]; then
    add_non_feature_branch "$branch_name"
  elif [ "$option" == "--remove" ]; then
    remove_non_feature_branch "$branch_name"
  else
    echo_inline_error "unsupported option '$option'"
    echo_non_feature_branch_usage
    exit 1
  fi
}


# Returns git-town configuration
function get_configuration {
  local config_setting_name=$1
  git config "git-town.$config_setting_name"
}


# Remove a non-feature branch if possible, and show confirmation
function remove_non_feature_branch {
  local branch_name=$1

  if [ -z "$branch_name" ]; then
    echo_inline_error "missing branch name"
    echo_non_feature_branch_usage
    exit 1
  elif [ "$(is_non_feature_branch "$branch_name")" = false ]; then
    echo_inline_error "'$branch_name' is not a non-feature branch"
    exit 1
  else
    local new_branches=$(remove_string "$non_feature_branch_names" ',' "$branch_name")
    store_configuration non-feature-branch-names "$new_branches"
  fi
}


function show_config {
  echo_inline_bold "Main branch: "
  show_main_branch
  echo_inline_bold "Non-feature branches:"
  if [ -n "$non_feature_branch_names" ]; then
    echo
    split_string "$non_feature_branch_names" ","
  else
    echo ' [none]'
  fi
}


function show_main_branch {
  if [ -n "$main_branch_name" ]; then
    echo "$main_branch_name"
  else
    echo '[none]'
  fi
}


function show_non_feature_branches {
  if [ -n "$non_feature_branch_names" ]; then
    split_string "$non_feature_branch_names" ","
  fi
}


function show_or_update_main_branch {
  local branch_name=$1
  if [ -n "$branch_name" ]; then
    if [ "$(has_branch "$branch_name")" = true ]; then
      store_configuration main-branch-name "$branch_name"
    else
      echo_inline_error "no branch named '$branch_name'"
      exit 1
    fi
  else
    show_main_branch
  fi
}


function show_or_update_non_feature_branches {
  local operation=$1
  local branch_name=$2
  if [ -n "$operation" ]; then
    add_or_remove_non_feature_branches "$operation" "$branch_name"
  else
    show_non_feature_branches
  fi
}


# Persists git-town configuration
function store_configuration {
  local config_setting_name=$1
  local value=$2
  git config "git-town.$config_setting_name" "$value"

  # update $main_branch_name and $non_feature_branch_names accordingly
  if [ $? == 0 ]; then
    if [ "$config_setting_name" == "main-branch-name" ]; then
      main_branch_name="$value"
    elif [ "$config_setting_name" == "non-feature-branch-names" ]; then
      non_feature_branch_names="$value"
    fi
  fi
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
