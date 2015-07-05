#!/usr/bin/env bash


# Helper methods for managing the configuration of which branches
# are cut from which ones


# Returns the names of all branches that are registered in the hierarchy metadata,
# as an iterable list
function all_registered_branches {
  git config --get-regexp "git-town\.branches\.parent" | cut -d ' ' -f 1 | sed 's/^git-town.branches.parent\.//' | sort | uniq
}


# Returns the names of all ancestor branches of the given branch
function ancestor_branches {
  local branch_name=$1

  ancestors=$(git config "git-town.branches.ancestors.$branch_name")
  escaped=${ancestors/ /, /}
  echo "$escaped"
}


# Returns the names of all branches that have this branch as their immediate parent
function child_branches {
  local branch_name=$1
  git config --get-regexp "^git-town\.branches\.parent\." | grep "$branch_name$" | sed 's/git-town\.branches\.parent\.//' | sed "s/ $branch_name$//"
}


# Calculates the "ancestors" property for the given branch
# out of the existing "parent" properties
function compile_ancestor_branches {
  local current_branch=$1

  # delete the existing entry
  delete_ancestors_entry "$current_branch"

  # re-create it from scratch
  local ancestor_branches=''
  local parent
  while [ "$current_branch" != "$MAIN_BRANCH_NAME" ]; do
    parent=$(parent_branch "$current_branch")
    ancestor_branches="$parent $ancestor_branches"
    current_branch=$parent
  done

  # truncate the trailing comma
  # shellcheck disable=SC2001
  ancestor_branches=$(echo "$ancestor_branches" | sed 's/ $//')

  # save the result into the configuration
  git config git-town.branches.ancestors."$(normalized_branch_name "$1")" "$ancestor_branches"
}


# Removes all ancestor cache entries
function delete_all_ancestor_entries {
  git config --remove-section git-town.branches.ancestors
}


# Removes the "parent" entry for the given branch from the configuration
function delete_parent_entry {
  local branch_name=$1

  local normalized_branch ; normalized_branch=$(normalized_branch_name "$branch_name")
  if [ "$(knows_parent_branch "$normalized_branch")" == "true" ]; then
    git config --unset "git-town.branches.parent.$normalized_branch"
  fi

  delete_ancestors_entry "$branch_name"
}


# Removes the "ancestors" entry from the configuration
function delete_ancestors_entry {
  local branch_name=$1

  local normalized_branch ; normalized_branch=$(normalized_branch_name "$branch_name")
  if [ "$(knows_all_ancestor_branches "$normalized_branch")" == "true" ]; then
    git config --unset "git-town.branches.ancestors.$normalized_branch"
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

  if [ "$(knows_all_ancestor_branches "$current_branch")" = false ]; then
    # Here we don't have the parent branches list --> make sure we know all ancestors, then recompile it from all ancestors
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
          branch_numbers[number]="$branch_name"
        }

        print_branch 1 "$MAIN_BRANCH_NAME"
        i=2
        for branch in $branches; do
          if [ "$branch" != "$current_branch" -a "$branch" != "$MAIN_BRANCH_NAME" ]; then
            print_branch $i "$branch"
            i=$(( i + 1 ))
          fi
        done

        echo
        echo -n "Please enter the parent branch name or number for $(echo_inline_cyan_bold "$current_branch") ($MAIN_BRANCH_NAME): "
        read parent
        re='^[0-9]+$'
        if [[ $parent =~ $re ]] ; then
          # user entered a number here
          parent=${branch_numbers[$parent]}
        elif [ -z "$parent" ]; then
          # user entered nothing
          parent=$MAIN_BRANCH_NAME
        fi
        if [ "$(has_branch "$parent")" == "false" ]; then
          echo_error_header
          echo_error "branch '$parent' doesn't exist"
          exit_with_error newline
        fi
        store_parent_branch "$current_branch" "$parent"
        echo
      fi
      current_branch=$parent
    done
    compile_ancestor_branches "$1"
  fi
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
  if [ -z "$(git config --get git-town.branches.parent."$branch_name")" ]; then
    echo false
  else
    echo true
  fi
}


# Returns whether we know the parent branches for the given branch
function knows_all_ancestor_branches {
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
function oldest_ancestor_branch {
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
