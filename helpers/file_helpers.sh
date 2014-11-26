#!/bin/bash

# Helper methods for dealing with files and temp files.


# Unique string that identifies the current directory and git branch
temp_filename_suffix="$(pwd | tr '/' '_')"

# Path to the temp file used by these scripts.
user_input_filename="/tmp/git-town-user-input$temp_filename_suffix"


# Returns the path to the abort script for the given command
function abort_script_filename_for_command {
  script_filename "$1" 'abort'
}


# Returns the path to the continue script for the given command
function continue_script_filename_for_command  {
  script_filename "$1" 'continue'
}


# Removes the temp file.
function delete_temp_file {
  rm "$user_input_filename"
}


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


function script_filename {
  echo "/tmp/git_${1//-/_}_$2_$temp_filename_suffix"
}


# Returns the path to the continue script for the given command
function skip_script_filename_for_command  {
  script_filename "$1" 'skip'
}


# Returns the path to the undo script for the given command
function undo_script_filename_for_command {
  script_filename "$1" 'undo'
}
