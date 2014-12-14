#!/bin/bash

# Helper methods for dealing with files and temp files.


# Unique string that identifies the current directory
temp_filename_suffix="$(pwd | tr '/' '_')"


# Command lists
export steps="/tmp/${program}_${temp_filename_suffix}"
export undo_steps="/tmp/${program}_undo_${temp_filename_suffix}"


function add_step {
  append_to_file "$1" "$steps"
}


function add_undo_step {
  prepend_to_file "$1" "$undo_steps"
}


function append_to_file {
  local content=${1//\'/\"}
  local file=$2
  local operator=">"
  if [ -e "$file" ]; then operator=">>"; fi
  eval "echo '$content' $operator $file"
}


function prepend_to_file {
  local content=$1
  local file=$2
  if [ "$(has_file "$file")" = true ]; then
    local temp=$(temp_filename)
    echo "$1" | cat - "$file" > "$temp" && mv "$temp" "$file"
  else
    append_to_file "$@"
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
  echo "/tmp/git-town$RANDOM.tmp"
}
