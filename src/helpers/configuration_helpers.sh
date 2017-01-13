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


# Add or remove the default alias for the given command
function add_or_remove_alias {
  local command=$1
  local boolean=$2
  local previousAlias="$(git config --global "alias.$command")"

  if [ "$boolean" = true ]; then
    run_command "git config --global alias.$command town-$command"
  else
    local previousAlias="$(git config --global "alias.$command")"
    if [ "$previousAlias" = "town-$command" ]; then
      run_command "git config --global --unset alias.$command"
    fi
  fi
}


# Add or remove the default alias for all git town commands
function add_or_remove_aliases {
  local boolean=$1

  if [ "$boolean" != 'true' ] && [ "$boolean" != 'false' ]; then
    echo "Invalid alias boolean: '$boolean'."
    echo "Valid alias booleans are 'true' and 'false'."
  fi

  declare -a commands=('append' 'hack' 'kill' 'new-pull-request' 'prepend' 'prune-branches' 'rename-branch' 'repo' 'ship' 'sync')
  for command in "${commands[@]}"; do
    add_or_remove_alias "$command" "$boolean"
  done

  echo # trailing newline in case any commands run
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


function show_branch_tree {
  local branch=$1
  local indentation=$2

  for (( i=0; i < "$indentation"; i++)); do
    printf '  '
  done
  echo "$branch"
  local child_indentation=$(( indentation + 1 ))
  for child in $(child_branches "$branch"); do
    show_branch_tree "$child" "$child_indentation"
  done
}


function show_config {
  echo_bold_underline "Main branch:"
  echo_indented "$(show_main_branch)"
  echo

  echo_bold_underline "Perennial branches:"
  echo_indented "$(show_perennial_branches)"
  echo

  if [ -n "$MAIN_BRANCH_NAME" ]; then
    echo_bold_underline "Branch Ancestry:"
    show_branch_tree "$MAIN_BRANCH_NAME" 0 | indent
    echo
  fi

  echo_bold_underline "Pull branch strategy:"
  # shellcheck disable=SC2119
  echo_indented "$(show_or_update_pull_branch_strategy)"
  echo

  echo_bold_underline "git-hack push flag:"
  # shellcheck disable=SC2119
  echo_indented "$(show_or_update_hack_push_flag)"
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
  else
    echo '[none]'
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


# shellcheck disable=SC2120
function show_or_update_hack_push_flag {
  local flag=$1
  if [ -z "$flag" ]; then
    echo "$HACK_PUSH_FLAG"
  elif [ "$flag" != 'true' ] && [ "$flag" != 'false' ]; then
    echo "Invalid git-hack push flag: '$flag'."
    echo "Valid git-hack push flags are 'true' and 'false'."
  else
    store_configuration hack-push-flag "$flag"
  fi
}


# Update the pull branch strategy branch if a strategy is specified,
# otherwise show the current pull branch strategy
# shellcheck disable=SC2120
function show_or_update_pull_branch_strategy {
  local strategy=$1
  if [ -z "$strategy" ]; then
    echo "$PULL_BRANCH_STRATEGY"
  elif [ "$strategy" != 'merge' ] && [ "$strategy" != 'rebase' ] ; then
    echo "Invalid pull branch strategy: '$strategy'."
    echo "Valid pull branch strategies are 'merge' and 'rebase'."
  else
    store_configuration pull-branch-strategy "$strategy"
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

  # update $MAIN_BRANCH_NAME and $PERENNIAL_BRANCH_NAMES accordingly
  if git config "git-town.$config_setting_name" "$value"; then
    if [ "$config_setting_name" == "main-branch-name" ]; then
      MAIN_BRANCH_NAME="$value"
    elif [ "$config_setting_name" == "perennial-branch-names" ]; then
      PERENNIAL_BRANCH_NAMES="$value"
    fi
  fi
}


function undo_steps_for_add_perennial_branch {
  local branch_name=$1
  echo "remove_perennial_branch $branch_name"
}


function undo_steps_for_remove_perennial_branch {
  local branch_name=$1
  echo "add_perennial_branch $branch_name"
}


function update_parent_branch {
  local child_branch_name=$1
  local parent_branch_name=$2
  if [ "$(has_branch "$child_branch_name")" = false ]; then
    echo_inline_error "no branch named '$child_branch_name'"
    exit_with_error
  elif [ "$(has_branch "$parent_branch_name")" = false ]; then
    echo_inline_error "no branch named '$parent_branch_name'"
    exit_with_error
  fi
  delete_ancestors_entry "$child_branch_name"
  store_parent_branch "$child_branch_name" "$parent_branch_name"
}
