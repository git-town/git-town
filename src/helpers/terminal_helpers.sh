#!/usr/bin/env bash

# Helper methods for writing to the terminal.


function prompt_yn () {
  while true; do
    read yn
    case $yn in
      [Yy] ) return 0;;
      [Nn] ) return 1;;
      *    ) echo "Please answer yes (y) or no (n).";;
    esac
  done
}


# Prints a line in bold
function echo_bold {
  output_style_bold
  echo "$*"
  output_style_reset
}


# Prints a header line into the terminal.
function echo_header {
  echo
  echo_bold "$*"
}


# Prints an error header into the terminal.
function echo_error_header {
  local str=$(echo_indented Error)
  echo
  echo_red_bold "$str"
}


# Prints the provided error message
function echo_error {
  local str=$(echo_indented "$*")
  echo_red "$str"
}


# Prints the message idented
function echo_indented {
  echo "  $*"
}

# Prints an inline usage
function echo_inline_bold {
  output_style_bold
  printf "%s" "$*"
  output_style_reset
}


# Prints an inline error
function echo_inline_error {
  echo_red "error: $*"
}


# Prints an inline usage
function echo_inline_usage {
  echo "usage: $*"
}


# Prints a continutation of an inline usage
function echo_inline_usage_or {
  echo "   or: $*"
}


# Outputs the given text in red
function echo_red {
  output_style_red
  echo "$*"
  output_style_reset
}


# Outputs the given text in red and bold
function echo_red_bold {
  output_style_bold
  output_style_red
  echo "$*"
  output_style_reset
}


# Prints the provided usage message
function echo_usage {
  echo_indented "$*"
}


# Prints the header for the usage help screen into the terminal.
function echo_usage_header {
  local str=$(echo_indented Usage)
  echo_header "$str"
}


function exit_with_error {
  exit_with_status 1 "$1"
}


function exit_with_abort {
  exit_with_status 2 "$1"
}


# Exits the currently running script with an exit code.
function exit_with_status {
  if [ "$2" = "newline" ]; then
    echo
  fi
  exit "$1"
}


function output_style_bold {
  tput bold
}


function output_style_red {
  tput setaf 1
}


function output_style_reset {
  tput sgr0
}


# Prints a command
function print_command {
  local branch_name=$(get_current_branch_name)
  echo_header "[$branch_name] $*"
}


# Run a command, prints command and output
function run_command {
  local cmd="$*"
  print_command "$cmd"
  eval "$cmd" 2>&1
}
