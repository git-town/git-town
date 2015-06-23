#!/usr/bin/env bash


# Helper methods for managing the configuration of which branches
# are cut from which ones


# Returns the names of all branches that are registered in the hierarchy metadata,
# as an iterable list
function all_registered_branches {
  git config --get-regexp "git-town\.branches\.parent" | cut -d ' ' -f 1 | sed 's/^git-town.branches.parent\.//' | sort | uniq
}



# Returns the names of all branches that have this branch as their immediate parent
function child_branches {
  local current_branch=$1
  git config --get-regexp "^git-town\.branches\.parent\." | grep "$current_branch$" | sed 's/git-town\.branches\.parent\.//' | sed "s/ $current_branch$//"
}


# Calculates the "ancestors" property for the given branch
# out of the existing "parent" properties
function compile_ancestor_branches {
  local current_branch=$1

  # delete the existing entry
  delete_ancestors_entry "$current_branch"

  # re-create it from scratch
  local all_parent_branches=''
  local parent
  while [ "$current_branch" != "$MAIN_BRANCH_NAME" ]; do
    parent=$(parent_branch "$current_branch")
    all_parent_branches="$parent $all_parent_branches"
    current_branch=$parent
  done

  # truncate the trailing comma
  # shellcheck disable=SC2001
  all_parent_branches=$(echo "$all_parent_branches" | sed 's/ $//')

  # save the result into the configuration
  git config git-town.branches.ancestors."$(normalized_branch_name "$1")" "$all_parent_branches"
}


# Removes the "parent" entry from the configuration
function delete_parent_entry {
  local branch_name=$1

  local normalized_branch ; normalized_branch=$(normalized_branch_name "$branch_name")
  if [ "$(knows_parent_branch "$normalized_branch")" == "true" ]; then
    git config --unset "git-town.branches.parent.$normalized_branch"
  fi
}


# Removes the "ancestors" entry from the configuration
function delete_ancestors_entry {
  local branch_name=$1

  local normalized_branch ; normalized_branch=$(normalized_branch_name "$branch_name")
  if [ "$(knows_all_parent_branches "$normalized_branch")" == "true" ]; then
    git config --unset "git-town.branches.ancestors.$normalized_branch"
  fi
}


# Makes sure that we know all the parent branches
# Asks the user if necessary
# Aborts the script if not all branches become known.
function ensure_knows_parent_branches {
  local current_branch=$1

  if [ "$(knows_all_parent_branches "$current_branch")" = false ]; then
    # Here we don't have the parent branches list --> make sure we know all ancestors, then recompile it from all ancestors
    local parent

    while [ "$current_branch" != "$MAIN_BRANCH_NAME" ]; do
      if [ "$(knows_parent_branch "$current_branch")" = true ]; then
        parent=$(parent_branch "$current_branch")
      else
        # here we don't know the parent of the current branch -> ask the user
        echo
        echo -n "Please enter the parent branch for $(echo_inline_cyan_bold "$current_branch") ($(echo_inline_dim "$MAIN_BRANCH_NAME")): "
        read parent
        if [ -z "$parent" ]; then
          parent=$MAIN_BRANCH_NAME
        fi
        if [ "$(has_branch "$parent")" == "false" ]; then
          echo_error_header
          echo_error "branch '$parent' doesn't exist"
          exit_with_error newline
        fi
        store_parent_branch "$current_branch" "$parent"
      fi
      current_branch=$parent
    done
    compile_ancestor_branches "$1"
  fi
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
  if [ -z "$(git config --get git-town.branches.ancestors."$branch_name")" ]; then
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


# Returns the name of the branch from the branch hierarchy
# that is the direct ancestor of main
function oldest_parent_branch {
  git config --get "git-town.branches.ancestors.$(normalized_branch_name "$branch_name")" | cut -d ' ' -f 2
}


# Returns the names of all parent branches, in hierarchical order
function parent_branch {
  local branch_name=$1
  git config --get "git-town.branches.parent.$(normalized_branch_name "$branch_name")"
}


# Returns the names of all parent branches of the given branch,
# as a string list, in hierarchical order,
function parent_branches {
  local branch_name=$1
  git config --get "git-town.branches.ancestors.$(normalized_branch_name "$branch_name")" | tr ' ' '\n'
}


# Stores the given branch as the parent branch for the given branch
function store_parent_branch {
  local branch=$1
  local parent_branch=$2
  git config "git-town.branches.parent.$(normalized_branch_name "$branch")" "$parent_branch"
}


function undo_steps_for_delete_parent_entry {
  local branch_name=$1

  if [ "$(knows_parent_branch "$branch_name")" == "true" ]; then
    echo "store_parent_branch $branch_name $(parent_branch "$branch_name")"
  fi
}


function undo_steps_for_store_parent_branch {
  local branch=$1

  local old_parent_branch ; old_parent_branch=$(parent_branch "$branch")
  echo "store_parent_branch $branch $old_parent_branch"
}
