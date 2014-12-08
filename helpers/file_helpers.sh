#!/bin/bash

# Helper methods for dealing with files and temp files.


# Unique string that identifies the current directory
temp_filename_suffix="$(pwd | tr '/' '_')"

# Command lists
export command_list="/tmp/${program}_${temp_filename_suffix}"
export abort_command_list="/tmp/${program}_abort_${temp_filename_suffix}"
export undo_command_list="/tmp/${program}_undo_${temp_filename_suffix}"


function add_to_abort_command_list {
  prepend_to_file "$1" "$abort_command_list"
}


function add_to_command_list {
  add_to_file "$1" "$command_list"
}


function add_to_file {
  local content=$1
  local file=$2
  local operator=">"
  if [ -e "$file" ]; then operator=">>"; fi
  eval "echo '$content' $operator $file"
}


function add_to_undo_command_list {
  prepend_to_file "$1" "$undo_command_list"
}


function prepend_to_file {
  local content=$1
  local file=$2
  if [ "$(has_file "$file")" = true ]; then
    local temp=$(temp_filename)
    echo "$1" | cat - "$file" > "$temp" && mv "$temp" "$file"
  else
    add_to_file "$@"
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
