#!/usr/bin/env bash


# Creates a new branch with the given name off the given parent branch
function create_branch {
  local new_branch_name=$1
  local parent_branch_name=$2
  run_command "git branch $new_branch_name $parent_branch_name"
}


# Creates and checkouts a new branch with the given name off the given parent branch
function create_and_checkout_branch {
  local new_branch_name=$1
  local parent_branch_name=$2
  run_command "git checkout -b $new_branch_name $parent_branch_name"
}


# Deletes the local branch with the given name
function delete_local_branch {
  local branch_name=$1
  local op="d"
  if [ "$2" == "force" ]; then op="D"; fi
  run_command "git branch -$op $branch_name"
}


# Deletes the remote branch with the given name
function delete_remote_branch {
  local branch_name=$1
  run_command "git push origin :${branch_name}"
}


# Exits if the repository has a branch with the given name
function ensure_does_not_have_branch {
  local branch_name=$1
  if [ "$(has_branch "$branch_name")" = true ]; then
    echo_error_header
    echo_error "A branch named '$branch_name' already exists"
    exit_with_error newline
  fi
}



# Exits if the repository does not have a branch with the given name
function ensure_has_branch {
  local branch_name=$1
  if [ "$(has_branch "$branch_name")" == false ]; then
    echo_error_header
    echo_error "There is no branch named '$branch_name'"
    exit_with_error newline
  fi
}


# Exits if any of the branches do not exist
function ensure_has_branches {
  local branches=$1

  split_string "$branches" ',' | while read branch; do
    ensure_has_branch "$branch"
  done
}


# Returns the current branch name
function get_current_branch_name {
  if [ "$(rebase_in_progress)" = true ]; then
    sed 's/^refs\/heads\///' < .git/rebase-apply/head-name
  else
    git rev-parse --abbrev-ref HEAD
  fi
}


# Returns true if the repository has a branch with the given name
function has_branch {
  local branch_name=$1
  if [ "$(git branch -a | tr -d '* ' | sed 's/remotes\/origin\///' | grep -c "^$branch_name\$")" = 0 ]; then
    echo false
  else
    echo true
  fi
}


# Returns the names of local branches
function local_branches {
  git branch | tr -d ' ' | sed 's/\*//g'
}


# Returns the names of local branches without the main branch
function local_branches_without_main {
  local_branches | grep -v "^$main_branch_name\$"
}


# Returns the names of local branches that have been merged into main
function local_merged_branches {
  git branch --merged "$main_branch_name" | tr -d ' ' | sed 's/\*//g'
}


# Pushes the branch with the given name to origin
function push_branch {
  local branch_name=$1
  local force=$2
  if [ "$(has_tracking_branch "$branch_name")" = true ]; then
    if [ "$(needs_push "$branch_name")" = true ]; then
      if [ -n "$force" ]; then
        run_command "git push -f origin $branch_name"
      else
        run_command "git push origin $branch_name"
      fi
    fi
  else
    run_command "git push -u origin $branch_name"
  fi
}


# Returns the names of remote branches that have been merged into main
function remote_merged_branches {
  git branch -r --merged "$main_branch_name" | grep -v HEAD | tr -d ' ' | sed 's/origin\///g'
}


# Returns the names of remote branches that have been merged into main
# that have not been checked out locally
function remote_only_merged_branches {
  local local_temp=$(temp_filename)
  local remote_temp=$(temp_filename)
  local_merged_branches > "$local_temp"
  remote_merged_branches > "$remote_temp"
  comm -13 <(sort "$local_temp") <(sort "$remote_temp")
  rm "$local_temp"
  rm "$remote_temp"
}


function undo_steps_for_create_and_checkout_feature_branch {
  local branch=$(get_current_branch_name)
  local branch_to_create="$1"
  echo "checkout $branch"
  echo "delete_local_branch $branch_to_create"
}


function undo_steps_for_delete_local_branch {
  local branch_to_delete="$1"
  local sha=$(sha_of_branch "$branch_to_delete")
  echo "create_branch $branch_to_delete $sha"
}


function undo_steps_for_delete_remote_branch {
  local branch_to_delete="$1"
  echo "push_branch $branch_to_delete"
}

