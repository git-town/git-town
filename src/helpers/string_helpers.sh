#!/bin/bash

# Helper functions for working with strings


# Split string with delimiter onto separate lines
function split_string {
  local string=$1
  local delimiter=$2

  echo "$string" | tr "$delimiter" '\n' | sed 's/^ *//;s/ *$//'
}


# Join a previously split string with delimiter
function join_string {
  local string=$1
  local delimiter=$2

  echo "$string" | tr '\n' "$delimiter" | sed "s/^ *$delimiter//;s/$delimiter *$//"
}


# Insert string into a delimiter-separated string
function insert_string {
  local delimiter_separated_string=$1
  local delimiter=$2
  local string_to_insert=$3

  local split="$(split_string "$delimiter_separated_string" "$delimiter")"
  join_string "$split"$'\n'"$string_to_insert" "$delimiter"
}


# Remove string from delimiter-separated string
function remove_string {
  local delimiter_separated_string=$1
  local delimiter=$2
  local string_to_remove=$3

  local split="$(split_string "$delimiter_separated_string" "$delimiter")"
  join_string "$(echo "$split" | sed "s/^${string_to_remove}$//;/^$/d")" "$delimiter"
}


function parameters_as_string {
  local str=""
  for arg in "$@"; do
    if [ "$arg" != "${arg/ /}" ]; then
      # Wrap arg in double quotes, escape double quotes within arg
      arg="\"${arg//\"/\\\"}\""
    fi
    str="$str $arg"
  done
  echo "${str/ /}" # Remove initial space
}
