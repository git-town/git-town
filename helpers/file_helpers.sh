#!/bin/bash

# Helper methods for dealing with files and temp files.


# Unique string that identifies the current directory
temp_filename_suffix="$(pwd | tr '/' '_')"

# Temporary filename used for short term storage of user input
export user_input_filename="/tmp/git-town-user-input_${temp_filename_suffix}"

# Scripts filenames
for action in "abort" "continue" "undo"; do
  declare -r ${action}_script_filename="/tmp/${program}_${action}_${temp_filename_suffix}"
done
