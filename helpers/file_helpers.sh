#!/bin/bash

# Helper methods for dealing with files and temp files.


# Unique string that identifies the current directory
temp_filename_suffix="$(pwd | tr '/' '_')"

# Temporary filename used for short term storage of user input
export user_input_filename="/tmp/git-town-user-input_${temp_filename_suffix}"

# Scripts filenames
actions=(abort continue undo)
for action in "${actions[@]}"; do
  declare ${action}_script_filename="/tmp/${program}_${action}_${temp_filename_suffix}"
done


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
