#!/usr/bin/env bash

# Helper methods for dealing with configuration.

function echo_non_feature_branch_usage {
  echo_inline_usage 'git town non-feature-branches (--add | --remove) <branchname>'
}


# Add a new non-feature branch if possible
function add_non_feature_branch {
  local branch_name=$1

  if [ -z "$branch_name" ]; then
    echo_inline_error "missing branch name"
    echo_non_feature_branch_usage
    exit_with_error
  elif [ "$(has_branch "$branch_name")" = false ]; then
    echo_inline_error "no branch named '$branch_name'"
    exit_with_error
  elif [ "$(is_non_feature_branch "$branch_name")" = true ]; then
    echo_inline_error "'$branch_name' is already a non-feature branch"
    exit_with_error
  elif [ "$branch_name" == "$MAIN_BRANCH_NAME" ]; then
    echo_inline_error "'$branch_name' is already set as the main branch"
    exit_with_error
  else
    local new_branches=$(insert_string "$NON_FEATURE_BRANCH_NAMES" ',' "$branch_name")
    store_configuration non-feature-branch-names "$new_branches"
  fi
}


# Add or remove non-feature branch if possible
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
    exit_with_error
  fi
}


# Ensure that non-feature branches don't contain main branch
function ensure_valid_non_feature_branches {
  local branches=$1

  split_string "$branches" ',' | while read branch; do
    if [[ "$branch" == "$MAIN_BRANCH_NAME" ]]; then
      echo_error_header
      echo_error "'$branch' is already set as the main branch"
      exit_with_error
    fi
  done
}


# Returns git-town configuration
function get_configuration {
  local config_setting_name=$1
  git config "git-town.$config_setting_name"
}


# Returns whether or not Git Town is configured
function is_git_town_configured {
  if [ -n "$MAIN_BRANCH_NAME" ] && get_configuration 'non-feature-branch-names'; then
    echo true
  else
    echo false
  fi
}


# Remove all Git Town configuration
function remove_all_configuration {
  # output is redirected because git will print an error
  # to stdout if the config section doesn't exist
  git config --remove-section git-town > /dev/null 2>&1
}


# Remove a non-feature branch if possible
function remove_non_feature_branch {
  local branch_name=$1

  if [ -z "$branch_name" ]; then
    echo_inline_error "missing branch name"
    echo_non_feature_branch_usage
    exit_with_error
  elif [ "$(is_non_feature_branch "$branch_name")" = false ]; then
    echo_inline_error "'$branch_name' is not a non-feature branch"
    exit_with_error
  else
    local new_branches=$(remove_string "$NON_FEATURE_BRANCH_NAMES" ',' "$branch_name")
    store_configuration non-feature-branch-names "$new_branches"
  fi
}


# Begin the git town setup wizard
function setup_configuration {
  setup_configuration_main_branch
  echo
  setup_configuration_non_feature_branches
  echo "Done with configuration:"
  show_config | indent
}


# Ask and store main-branch-name
function setup_configuration_main_branch {
  echo "Please specify the main dev branch (typically 'master' or 'development'):"
  read main_branch_input
  if [[ -z "$main_branch_input" ]]; then
    echo_error_header
    echo_error "You have not provided the name for the main branch."
    echo_error "Aborting Git Town configuration."
    exit_with_error newline
  fi

  ensure_has_branch "$main_branch_input" || exit_with_error
  store_configuration main-branch-name "$main_branch_input"
}


# Ask and store non-feature-branch-names
function setup_configuration_non_feature_branches {
  echo "Git Town supports non-feature branches like 'release' or 'production'."
  echo "These branches cannot be shipped and will not merge '$MAIN_BRANCH_NAME' when syncing."
  echo "Please enter your non-feature branches as a comma separated list or a blank line to skip."
  echo "Example: 'qa, production'"
  read non_feature_input

  if [[ -n "$non_feature_input" ]]; then
    ensure_has_branches "$non_feature_input" || exit_with_error
    ensure_valid_non_feature_branches "$non_feature_input" || exit_with_error
  fi

  store_configuration non-feature-branch-names "$non_feature_input"
}


# Perform `git town config` operation ("reset", "setup", "show")
function run_config_operation {
  local operation=$1

  if [ -n "$operation" ]; then
    if [ "$operation" == "--setup" ]; then
      setup_configuration
    elif [ "$operation" == "--reset" ]; then
      remove_all_configuration
    else
      echo "usage: git town config [--reset | --setup]"
    fi
  else
    show_config
  fi
}


function show_config {
  echo_inline_bold "Main branch: "
  show_main_branch
  echo_inline_bold "Non-feature branches:"
  if [ -n "$NON_FEATURE_BRANCH_NAMES" ]; then
    echo
    split_string "$NON_FEATURE_BRANCH_NAMES" ","
  else
    echo ' [none]'
  fi
}


function show_main_branch {
  if [ -n "$MAIN_BRANCH_NAME" ]; then
    echo "$MAIN_BRANCH_NAME"
  else
    echo '[none]'
  fi
}


function show_non_feature_branches {
  if [ -n "$NON_FEATURE_BRANCH_NAMES" ]; then
    split_string "$NON_FEATURE_BRANCH_NAMES" ","
  fi
}


# Update the main branch if a branch name is specified,
# otherwise show the current main branch name
function show_or_update_main_branch {
  local branch_name=$1
  if [ -n "$branch_name" ]; then
    if [ "$(has_branch "$branch_name")" = true ]; then
      store_configuration main-branch-name "$branch_name"
    else
      echo_inline_error "no branch named '$branch_name'"
      exit_with_error
    fi
  else
    show_main_branch
  fi
}


# Update the non-feature branches if a branch name and
# operation is specified, otherwise show the current
# non-feature branch names
function show_or_update_non_feature_branches {
  local operation=$1
  local branch_name=$2
  if [ -n "$operation" ]; then
    add_or_remove_non_feature_branches "$operation" "$branch_name"
  else
    show_non_feature_branches
  fi
}


# Persists the given git-town configuration setting
#
# The configuration setting is provided as a name-value pair, and
# the respective MAIN_BRANCH_NAME or NON_FEATURE_BRANCH_NAMES
# shell variable is updated.
function store_configuration {
  local config_setting_name=$1
  local value=$2
  git config "git-town.$config_setting_name" "$value"

  # update $MAIN_BRANCH_NAME and $NON_FEATURE_BRANCH_NAMES accordingly
  if [ $? == 0 ]; then
    if [ "$config_setting_name" == "main-branch-name" ]; then
      MAIN_BRANCH_NAME="$value"
    elif [ "$config_setting_name" == "non-feature-branch-names" ]; then
      NON_FEATURE_BRANCH_NAMES="$value"
    fi
  fi
}
