#!/usr/bin/env bash

# Helper methods for dealing with files and temp files.


function has_lines {
  local file="$1"

  if [ "$(has_file "$file")" = true ] && [ "$(number_of_lines "$file")" -gt 0 ]; then
    echo true
  else
    echo false
  fi
}


function has_file {
  local file="$1"

  if [ -n "$file" ] && [ -f "$file" ]; then
    echo true
  else
    echo false
  fi
}


function number_of_lines {
  local file="$1"
  wc -l < "$file" | tr -d ' '
}


function peek_line {
  local file="$1"
  head -n 1 "$file"
}


function pop_line {
  local file="$1"
  peek_line "$file"
  remove_line "$file"
}


function remove_line {
  local file="$1"
  local temp=$(temp_filename)
  tail -n +2 "$file" > "$temp"
  mv "$temp" "$file"
}


function prepend_to_file {
  local content=$1
  local file=$2
  if [ "$(has_file "$file")" = true ]; then
    local temp=$(temp_filename)
    echo "$content" | cat - "$file" > "$temp" && mv "$temp" "$file"
  else
    echo "$content" > "$file"
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
