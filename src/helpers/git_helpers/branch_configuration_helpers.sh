#!/usr/bin/env bash


# Helper methods for managing the configuration of which branches
# are cut from which ones


# Calculates the "parents" property for the given branch
# out of the existing "parent" properties
function compile_parent_branches {
  local current_branch=$1

  # delete the existing entry
  delete_parents_entry "$current_branch"

  # re-create it from scratch
  local all_parent_branches=''
  local parent
  while [ "$current_branch" != "$MAIN_BRANCH_NAME" ]; do
    parent=$(parent_branch "$current_branch")
    all_parent_branches="$parent,$all_parent_branches"
    current_branch=$parent
  done

  # truncate the trailing comma
  # shellcheck disable=SC2001
  all_parent_branches=$(echo "$all_parent_branches" | sed 's/,$//')

  # save the result into the configuration
  git config git-town.branches.parents."$(normalized_branch_name "$1")" "$all_parent_branches"
}


# Removes the "parents" entry from the configuration
function delete_parents_entry {
  local branch_name=$1
  git config --unset git-town.branches.parents."$branch_name"
}


# Returns whether we know the parent branch for the given branch
function knows_parent_branch {
  local branch_name=$1
  if [ -z "$(git config --get git-town.branches.parent."$branch_name")" ]; then
    echo false
  else
    echo true
  fi
}


# Returns whether we know the parent branches for the given branch
function knows_all_parent_branches {
  local branch_name=$1
  if [ -z "$(git config --get git-town.branches.parents."$branch_name")" ]; then
    echo false
  else
    echo true
  fi
}


# Returns the given branch name normalized so that it is compatible
# with Git's command-line interface for configuration data
function normalized_branch_name {
  local branch_name=$1
  echo "$branch_name" | tr '_' '-'
}


# Returns the names of all parent branches, in hierarchical order
function parent_branch {
  local branch_name=$1
  git config --get git-town.branches.parent."$(normalized_branch_name "$branch_name")"
}


# Returns the names of all parent branches,
# as a string list, in hierarchical order,
function parent_branches {
  local branch_name=$1
  git config --get git-town.branches.parents."$(normalized_branch_name "$branch_name")" | tr ',' '\n'
}


# Stores the given branch as the parent branch for the given branch
function store_parent_branch {
  local branch=$1
  local parent_branch=$2
  git config git-town.branches.parent."$branch" "$parent_branch"
}

