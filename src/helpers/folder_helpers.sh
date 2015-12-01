#!/usr/bin/env bash

# Helper methods for dealing with directories.


# Changes into the given directory
function change_directory {
  local directory=$1
  local omit_if_doesnt_exist=$2
  if [ -d "$directory" ] || [ -z "$omit_if_doesnt_exist" ]; then
    run_command "cd $directory"
  fi
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

