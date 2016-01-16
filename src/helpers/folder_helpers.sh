#!/usr/bin/env bash

# Helper methods for dealing with directories.


# Changes into the given directory
function change_directory {
  local directory=$1
  run_command "cd $directory"
}

function change_directory_if_exists {
  local directory=$1
  if [ -d "$directory" ]; then
    change_directory "$directory"
  fi
}


function undo_steps_for_change_directory {
  echo "change_directory $(pwd)"
}


function undo_steps_for_change_directory_if_exists {
  undo_steps_for_change_directory
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

