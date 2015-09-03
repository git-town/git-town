#!/usr/bin/env bash


# Abort a cherry-pick
function abort_cherry_pick {
  run_command "git cherry-pick --abort"
}


# Cherry picks the SHAs into the current branch
function cherry_pick {
  local SHAs="$*"
  run_command "git cherry-pick $SHAs"
}


# Determines whether the current branch has a cherry-pick in progress
function cherry_pick_in_progress {
  if [ "$(git status | grep -c "You are currently cherry-picking")" == 1 ]; then
    echo true
  else
    echo false
  fi
}


# Continues cherry-pick if one is in progress
function continue_cherry_pick {
  if [ "$(has_open_changes)" == true ]; then
    run_command "git commit --no-edit"
  fi

  if [ "$(cherry_pick_in_progress)" == true ]; then
    run_command "git cherry-pick --continue"
  fi
}
