#!/bin/bash

# Helper methods for dealing with files and temp files.


# Unique string that identifies the current directory
temp_filename_suffix="$(pwd | tr '/' '_')"

# Scripts filenames
for action in "abort" "continue" "undo"; do
  declare -r ${action}_script_filename="/tmp/${program}_${action}_${temp_filename_suffix}"
done


function temp_filename {
  local file=$(temp_filename_unsafe)
  while [ -e "$file" ]; do
    file=$(temp_filename_unsafe)
  done
  echo "$file"
}


function temp_filename_unsafe {
  echo "/tmp/git-town$RANDOM$RANDOM.tmp"
}
