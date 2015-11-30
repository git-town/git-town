#!/usr/bin/env bash


# Returns whether the given branch is in sync with its tracking branch
function branch_needs_push {
  local branch_name=$1
  local tracking_branch_name="origin/$branch_name"
  if [ "$(git rev-list --left-right "$branch_name...$tracking_branch_name" | wc -l | tr -d ' ')" != 0 ]; then
    echo true
  else
    echo false
  fi
}


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
  store_parent_branch "$new_branch_name" "$parent_branch_name"
}


# Creates a new remote branch for the given local branch
function create_tracking_branch {
  local branch_name=$1
  run_command "git push -u origin $branch_name"
}


# Deletes the local branch with the given name
function delete_local_branch {
  local branch_name=$1
  local op="d"
  if [ "$2" == "force" ] || [ "$(delete_local_branch_needs_force "$branch_name")" = true ]; then
    op="D"
  fi
  run_command "git branch -$op $branch_name"
}


# Deletes the remote branch with the given name
function delete_remote_branch {
  local branch_name=$1
  run_command "git push origin :${branch_name}"
}


# Deletes the remote branch with the given name
function delete_remote_only_branch {
  delete_remote_branch "$@"
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


# Returns the branch that a user would expect when running `git checkout -`
function expected_previous_branch {
  local initial_previous_branch=$1
  local initial_branch=$2

  # previous branch still exists
  if [ "$(has_local_branch "$initial_previous_branch")" = true ]; then

    # current branch is unchanged
    if [ "$(get_current_branch_name)" = "$initial_branch" ]; then
      echo "$initial_previous_branch"

    # current branch is deleted
    elif [ "$(has_local_branch "$initial_branch")" = false ]; then
      echo "$initial_previous_branch"

    # current branch is new
    else
      if [ "$(has_local_branch "$initial_branch")" = true ]; then
        echo "$initial_branch"
      else
        echo "$MAIN_BRANCH_NAME"
      fi
    fi

  # previous branch is deleted
  else
    echo "$MAIN_BRANCH_NAME"
  fi
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


# Returns true if the repository has a local branch with the given name
function has_local_branch {
  local branch_name=$1
  if [ "$(local_branches | grep -c "^$branch_name\$")" = 0 ]; then
    echo false
  else
    echo true
  fi
}


# Returns the names of local branches
function local_branches {
  git branch | tr -d ' ' | sed 's/\*//g'
}


# Returns the names of local branches
function local_branches_with_main_first {
  if [ -n "$MAIN_BRANCH_NAME" ]; then
    echo "$MAIN_BRANCH_NAME"
  fi
  local_branches_without_main
}


# Returns the names of local branches without the main branch
function local_branches_without_main {
  local_branches | grep -v "^$MAIN_BRANCH_NAME\$"
}


# Returns the names of local branches that have been merged into main
function local_merged_branches {
  git branch --merged "$MAIN_BRANCH_NAME" | tr -d ' ' | sed 's/\*//g'
}


# Returns whether or not the force flag is needed to delete the given branch
function delete_local_branch_needs_force {
  local branch_name=$1
  if [ -n "$(git log "..$branch_name")" ]; then
    echo true
  else
    echo false
  fi
}


function preserve_checkout_history {
  local initial_previous_branch_name=$1
  local initial_branch_name=$2
  local desired_previous_branch=$(expected_previous_branch "$initial_previous_branch_name" "$initial_branch_name")
  set_previous_branch "$desired_previous_branch"
}


# Returns the previously checked out branch name
function previous_branch_name {
  # --verify and --quiet needed to suppress errors when previous branch is unresolvable
  git rev-parse --verify --quiet --abbrev-ref "@{-1}"
}


# Pushes the branch with the given name to origin
function push_branch {
  local branch_name=$1
  local force=$2
  if [ "$(branch_needs_push "$branch_name")" = true ]; then
    if [ -n "$force" ]; then
      run_command "git push -f origin $branch_name"
    else
      if [ "$(get_current_branch_name)" = "$branch_name" ]; then
        run_command "git push"
      else
        run_command "git push origin $branch_name"
      fi
    fi
  fi
}


# Returns the names of remote branches that have been merged into main
function remote_merged_branches {
  git branch -r --merged "$MAIN_BRANCH_NAME" | grep -v HEAD | tr -d ' ' | sed 's/origin\///g'
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


function set_previous_branch {
  local desired_previous_branch=$1

  # nothing to do, exit early
  if [ "$desired_previous_branch" = "$(previous_branch_name)" ]; then
    return
  fi

  local current_branch="$(get_current_branch_name)"

  checkout_silently "$desired_previous_branch"
  checkout_silently "$current_branch"
}


function undo_steps_for_create_and_checkout_branch {
  local current_branch=$(get_current_branch_name)
  local branch_to_create="$1"
  local current_branch=$(get_current_branch_name)

  echo "checkout $current_branch"
  echo "delete_local_branch $branch_to_create"
  echo "delete_parent_entry $branch_to_create"
  echo "delete_ancestors_entry $branch_to_create"
}


function undo_steps_for_create_tracking_branch {
  local branch_name=$1
  echo "delete_remote_branch $1"
}


function undo_steps_for_delete_local_branch {
  local branch_to_delete="$1"
  local sha=$(sha_of_branch "$branch_to_delete")
  echo "create_branch $branch_to_delete $sha"
}


function undo_steps_for_delete_remote_branch {
  local branch_to_delete="$1"
  echo "create_tracking_branch $branch_to_delete"
}


function undo_steps_for_delete_remote_only_branch {
  local branch_to_delete="$1"
  local remote_sha="$(git log origin/"$branch_to_delete" | head -1 | cut -d ' ' -f 2)"
  echo "run_command 'git push origin $remote_sha:refs/heads/$branch_to_delete'"
}


function undo_steps_for_push_branch {
  echo "skip_current_branch_steps $UNDO_STEPS_FILE"
}
