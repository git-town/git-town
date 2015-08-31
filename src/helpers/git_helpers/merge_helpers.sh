#!/usr/bin/env bash


# Abort a merge
function abort_merge {
  run_command "git merge --abort"
}


# Continues merge if one is in progress
function continue_merge {
  if [ "$(has_open_changes)" == true ]; then
    run_command "git commit --no-edit"
  fi
}


# Merges the given branch into the current branch
function merge {
  local branch_name=$1
  run_command "git merge --no-edit $branch_name"
}


# Squash merges the given branch into the current branch
function squash_merge {
  local branch_name=$1
  run_command "git merge --squash $branch_name"
}


function commit_squash_merge {
  local branch_name=$1
  shift
  local options=$(parameters_as_string "$@")
  if ! [[ options == *"--author"* ]]; then
    squash_commit_author=''
    get_squash_commit_author "$branch_name"
    if [ "$squash_commit_author" != "$(local_author)" ]; then
      options="--author=\"$squash_commit_author\" $options"
    fi
  fi
  sed -i -e 's/^/# /g' .git/SQUASH_MSG
  run_command "git commit $options"
  if [ $? != 0 ]; then error_empty_commit; fi
}


function undo_steps_for_merge {
  echo "reset_to_sha $(current_sha) hard"
}


function post_undo_steps_for_commit_squash_merge {
  local current_branch_name=$(get_current_branch_name)
  echo "checkout $current_branch_name"
  echo "revert_commit $(git log -n 1 --format="%H")"
  if [ "$HAS_REMOTE" == true ]; then
    echo "push_branch $current_branch_name"
  fi
}
