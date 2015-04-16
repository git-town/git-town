#!/usr/bin/env bash

# Helper methods for dealing with directories.


# Changes into the given directory
function change_directory {
  local directory=$1
  run_git_command "cd $directory"
}


function undo_steps_for_change_directory {
  echo "change_directory $(pwd)"
}


# Returns whether the current working directory
# is in a subdirectory of the current Git workspace.
function is_in_git_sub_directory {
  if [ "$(pwd)" != "$(git_root)" ]; then
    echo true
  else
    echo false
  fi
}


# Returns the name of the root directory of the current Git workspace.
function git_root {
  git rev-parse --show-toplevel
}

