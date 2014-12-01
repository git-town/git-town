#!/bin/bash

# Helper methods for dealing with files and temp files.


# Unique string that identifies the current directory
temp_filename_suffix="$(pwd | tr '/' '_')"

# Scripts filenames
for action in "abort" "undo"; do
  declare -r ${action}_script_filename="/tmp/${program}_${action}_${temp_filename_suffix}"
done

export command_list_filename="/tmp/${program}_command_list_${temp_filename_suffix}"


# Ensures that the given tool is installed.
function ensure_tool_installed {
  local toolname=$1
  if [ "$(which "$toolname" | wc -l)" == 0 ]; then
    echo_error_header
    echo_error "You need the '$toolname' tool in order to run tests."
    echo_error "Please install it using your package manager,"
    echo_error "or on OS X with 'brew install $toolname'."
    exit_with_error
  fi
}


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
