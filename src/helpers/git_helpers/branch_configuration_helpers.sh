#!/usr/bin/env bash


# Helper methods for managing the configuration of which branches
# are cut from which ones


# Returns the names of all branches that are registered in the hierarchy metadata,
# as an iterable list
function all_registered_branches {
  git config --get-regexp "^git-town-branch\..*\.parent$" | cut -d ' ' -f 1 | sed 's/^git-town-branch\.\(.*\)\.parent$/\1/' | sort | uniq
}


# Returns the names of all parent branches of the given branch,
# as a space delimited string, in hierarchical order,
function ancestor_branches {
  local branch_name=$1
  git config --get "git-town-branch.$branch_name.ancestors"
}


# Returns the names of all branches that have this branch as their immediate parent
function child_branches {
  local branch_name=$1
  git config --get-regexp "^git-town-branch\..*\.parent$" | grep "$branch_name$" | cut -d ' ' -f 1 | sed 's/^git-town-branch\.\(.*\)\.parent$/\1/'
}


# Calculates the "ancestors" property for the given branch
# out of the existing "parent" properties
function compile_ancestor_branches {
  local branch_name=$1

  # delete the existing entry
  delete_ancestors_entry "$branch_name"

  # re-create it from scratch
  local ancestors=''
  local current_branch="$branch_name"
  while [ "$current_branch" != "$MAIN_BRANCH_NAME" ]; do
    local parent=$(parent_branch "$current_branch")
    ancestors="$parent $ancestors"
    current_branch=$parent
  done

  # truncate the trailing space
  # shellcheck disable=SC2001
  ancestors=$(echo "$ancestors" | sed 's/ $//')

  # save the result into the configuration
  git config "git-town-branch.$branch_name.ancestors" "$ancestors"
}


# Removes all ancestor cache entries
function delete_all_ancestor_entries {
  git config --get-regexp "^git-town-branch.*ancestors$" | cut -d ' ' -f 1 | while read ancestor_entry; do
    git config --unset "$ancestor_entry"
  done
}


# Removes the "parent" entry for the given branch from the configuration
function delete_parent_entry {
  local branch_name=$1
  if [ "$(knows_parent_branch "$branch_name")" == "true" ]; then
    git config --unset "git-town-branch.$branch_name.parent"
  fi
  delete_ancestors_entry "$branch_name"
}


# Removes the "ancestors" entry from the configuration
function delete_ancestors_entry {
  local branch_name=$1
  if [ "$(knows_all_ancestor_branches "$branch_name")" == "true" ]; then
    git config --unset "git-town-branch.$branch_name.ancestors"
  fi
}


# Updates the child branches of the given branch to point to the other given branch
function echo_update_child_branches {
  local branch=$1
  local new_parent=$2

  child_branches "$branch" | while read branch_name; do
    echo delete_ancestors_entry "$branch_name"
    echo store_parent_branch "$branch_name" "$new_parent"
  done
}


# Makes sure that we know all the parent branches
# Asks the user if necessary
# Aborts the script if not all branches become known.
function ensure_knows_parent_branches {
  local current_branch=$1

  if [ "$(knows_all_ancestor_branches "$current_branch")" = true ]; then
    return
  fi

  # Here we don't have the ancestors list --> make sure we know all ancestors, then recompile it from all ancestors
  local parent

  while [ "$current_branch" != "$MAIN_BRANCH_NAME" ]; do
    if [ "$(knows_parent_branch "$current_branch")" = true ]; then
      parent=$(parent_branch "$current_branch")
    else
      # here we don't know the parent of the current branch -> ask the user
      echo
      local branches=$(git for-each-ref --sort=-committerdate refs/heads/ --format='%(refname:short)')

      function print_branch {
        local number=$1
        local branch_name=$2

        output_style_bold
        printf "%3s: " "$number"
        output_style_reset
        echo "$branch_name"
      }

      echo
      echo "Feature branches can be branched directly off "
      echo "$MAIN_BRANCH_NAME or from other feature branches."
      echo
      echo "The former allows to develop and ship features completely independent of each other."
      echo "The latter allows to build on top of currently unshipped features."
      echo

      local branch_numbers
      print_branch 1 "$MAIN_BRANCH_NAME"
      branch_numbers[1]=$MAIN_BRANCH_NAME
      i=2
      for branch in $branches; do
        if [ "$branch" != "$current_branch" ] && [ "$branch" != "$MAIN_BRANCH_NAME" ]; then
          branch_numbers[i]=$branch
          print_branch $i "$branch"
          i=$(( i + 1 ))
        fi
      done

      local has_branch=false
      while [ $has_branch == false ]; do
        echo
        echo -n "Please specify the parent branch of $(echo_inline_cyan_bold "$current_branch") by name or number (default: $MAIN_BRANCH_NAME): "
        read parent
        re='^[0-9]+$'
        if [[ $parent =~ $re ]] ; then
          # user entered a number here
          parent=${branch_numbers[$parent]}
          if [ -z "$parent" ]; then
            echo_error_header
            echo_error "Invalid branch number"
          else
            has_branch=true
          fi
        elif [ -z "$parent" ]; then
          # user entered nothing
          parent=$MAIN_BRANCH_NAME
          has_branch=true
        else
          if [ "$(has_branch "$parent")" == "false" ]; then
            echo_error_header
            echo_error "branch '$parent' doesn't exist"
          else
            has_branch=true
          fi
        fi
      done
      store_parent_branch "$current_branch" "$parent"
      echo
    fi
    current_branch=$parent
  done
  compile_ancestor_branches "$1"
}


# Returns whether the given branch has child branches
function has_child_branches {
  local branch_name=$1

  if [ "$(child_branches "$branch_name")" == "" ]; then
    echo false
  else
    echo true
  fi
}


# Returns whether we know the parent branch for the given branch
function knows_parent_branch {
  local branch_name=$1
  if [ -z "$(parent_branch "$branch_name")" ]; then
    echo false
  else
    echo true
  fi
}


# Returns whether we know the parent branches for the given branch
function knows_all_ancestor_branches {
  local branch_name=$1
  if [ -z "$(ancestor_branches "$branch_name")" ]; then
    echo false
  else
    echo true
  fi
}


# Returns the name of the branch from the branch hierarchy
# that is the direct ancestor of main
function oldest_ancestor_branch {
  local branch_name=$1
  ancestor_branches "$branch_name" | cut -d ' ' -f 2
}


# Returns the names of all parent branches, in hierarchical order
function parent_branch {
  local branch_name=$1
  git config --get "git-town-branch.$branch_name.parent"
}


# Stores the given branch as the parent branch for the given branch
function store_parent_branch {
  local branch=$1
  local parent_branch=$2
  git config "git-town-branch.$branch.parent" "$parent_branch"
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
