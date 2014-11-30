#!/bin/bash


# Determines whether the given branch is ahead of main
function is_ahead_of_main {
  local branch_name=$1
  if [ "$(git log --oneline "$main_branch_name..$branch_name" | wc -l | tr -d ' ')" == 0 ]; then
    echo false
  else
    echo true
  fi
}
