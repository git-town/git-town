#!/usr/bin/env bash

# Helper methods for dealing with configuration.

# Add a new perennial branch if possible
function add_perennial_branch {
  local branch_name=$1

  if [ -z "$branch_name" ]; then
    echo_inline_error "missing branch name"
    echo_perennial_branch_usage
    exit_with_error
  elif [ "$(has_branch "$branch_name")" = false ]; then
    echo_inline_error "no branch named '$branch_name'"
    exit_with_error
  elif [ "$(is_perennial_branch "$branch_name")" = true ]; then
    echo_inline_error "'$branch_name' is already a perennial branch"
    exit_with_error
  elif [ "$branch_name" == "$MAIN_BRANCH_NAME" ]; then
    echo_inline_error "'$branch_name' is already set as the main branch"
    exit_with_error
  else
    local new_branches=$(insert_string "$PERENNIAL_BRANCH_NAMES" ' ' "$branch_name")
    store_configuration perennial-branch-names "$new_branches"
  fi
}


# Add or remove perennial branch if possible
function add_or_remove_perennial_branches {
  local option=$1
  local branch_name=$2

  if [ "$option" == "--add" ]; then
    add_perennial_branch "$branch_name"
  elif [ "$option" == "--remove" ]; then
    remove_perennial_branch "$branch_name"
  else
    echo_inline_error "unsupported option '$option'"
    echo_perennial_branch_usage
    exit_with_error
  fi
}


# Returns whether or not the perennial branches are configured
function are_perennial_branches_configured {
  if get_configuration 'perennial-branch-names'; then
    echo true
  else
    echo false
  fi
}


function echo_perennial_branch_usage {
  echo_inline_usage 'git town perennial-branches (--add | --remove) <branch_name>'
}


# Returns git-town configuration
function get_configuration {
  local config_setting_name=$1
  git config "git-town.$config_setting_name"
}


# Returns whether or not the main branch is configured
function is_main_branch_configured {
  if [ -n "$MAIN_BRANCH_NAME" ]; then
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


# Remove a perennial branch if possible
function remove_perennial_branch {
  local branch_name=$1

  if [ -z "$branch_name" ]; then
    echo_inline_error "missing branch name"
    echo_perennial_branch_usage
    exit_with_error
  elif [ "$(is_perennial_branch "$branch_name")" = false ]; then
    echo_inline_error "'$branch_name' is not a perennial branch"
    exit_with_error
  else
    local new_branches=$(remove_string "$PERENNIAL_BRANCH_NAMES" ' ' "$branch_name")
    store_configuration perennial-branch-names "$new_branches"
  fi
}


# Perform `git town config` operation
function run_config_operation {
  local operation=$1

  if [ -n "$operation" ]; then
    if [ "$operation" == "--setup" ]; then
      ensure_knows_configuration
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
  echo_inline_bold "Perennial branches:"
  if [ -n "$PERENNIAL_BRANCH_NAMES" ]; then
    echo
    split_string "$PERENNIAL_BRANCH_NAMES" " "
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


function show_perennial_branches {
  if [ -n "$PERENNIAL_BRANCH_NAMES" ]; then
    split_string "$PERENNIAL_BRANCH_NAMES" " "
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


# Update the perennial branches if a branch name and
# operation is specified, otherwise show the current
# perennial branch names
function show_or_update_perennial_branches {
  local operation=$1
  local branch_name=$2
  if [ -n "$operation" ]; then
    add_or_remove_perennial_branches "$operation" "$branch_name"
  else
    show_perennial_branches
  fi
}


# Persists the given git-town configuration setting
#
# The configuration setting is provided as a name-value pair, and
# the respective MAIN_BRANCH_NAME or PERENNIAL_BRANCH_NAMES
# shell variable is updated.
function store_configuration {
  local config_setting_name=$1
  local value=$2
  git config "git-town.$config_setting_name" "$value"

  # update $MAIN_BRANCH_NAME and $PERENNIAL_BRANCH_NAMES accordingly
  if [ $? == 0 ]; then
    if [ "$config_setting_name" == "main-branch-name" ]; then
      MAIN_BRANCH_NAME="$value"
    elif [ "$config_setting_name" == "perennial-branch-names" ]; then
      PERENNIAL_BRANCH_NAMES="$value"
    fi
  fi
}
