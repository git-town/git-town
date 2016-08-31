#!/usr/bin/env bash

# Helper functions for working with strings


function contains_string {
  local haystack=$1
  # shellcheck disable=SC2034
  local needle=$2
  if [[ $haystack =~ .*$needle.* ]]; then
    echo true
  else
    echo false
  fi
}


# Insert string into a delimiter-separated string
function insert_string {
  local delimiter_separated_string=$1
  local delimiter=$2
  local string_to_insert=$3

  local split="$(split_string "$delimiter_separated_string" "$delimiter")"
  join_string "$split"$'\n'"$string_to_insert" "$delimiter"
}


# Join a previously split string with delimiter
function join_string {
  local string=$1
  local delimiter=$2

  echo "$string" | tr '\n' "$delimiter" | sed "s/^ *$delimiter//;s/$delimiter *$//"
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


# Remove string from delimiter-separated string
function remove_string {
  local delimiter_separated_string=$1
  local delimiter=$2
  local string_to_remove=$3

  local split="$(split_string "$delimiter_separated_string" "$delimiter")"
  join_string "$(echo "$split" | sed "s/^${string_to_remove}$//;/^$/d")" "$delimiter"
}


# Split string with delimiter onto separate lines
function split_string {
  local string=$1
  local delimiter=$2

  echo "$string" | tr "$delimiter" '\n' | sed 's/^ *//;s/ *$//'
}


