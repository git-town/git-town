#!/usr/bin/env bash

shopt -s extglob

# Helper methods for writing to the terminal.


# Prints a line in bold
function echo_bold {
  output_style_bold
  echo "$@"
  output_style_reset
}


function echo_n_cyan_bold {
  output_style_cyan
  output_style_bold
  echo -n "$@"
  output_style_reset
}


function echo_n_dim {
  output_style_dim
  echo -n "$@"
  output_style_reset
}


# Prints a header line into the terminal.
function echo_header {
  echo
  echo_bold "$@"
}


# Prints an error header into the terminal.
function echo_error_header {
  echo
  echo_red_bold "$(echo_indented 'Error')"
}


# Prints the provided error message
function echo_error {
  echo_red "$(echo_indented "$@")"
}


# Prints the string if the condition is true
function echo_if_true {
  local string="$1"
  local condition="$2"

  if [ "$condition" = true ]; then
    echo "$string"
  fi
}


# Prints the message indented
function echo_indented {
  echo "$@" | indent
}


# Prints an inline usage
function echo_inline_bold {
  output_style_bold
  printf "%s" "$@"
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


# Prints a continuation of an inline usage
function echo_inline_usage_or {
  echo "or: $*" | indent 3
}


# Outputs the given text in red
function echo_red {
  output_style_red
  echo "$@"
  output_style_reset
}


# Outputs the given text in red and bold
function echo_red_bold {
  output_style_bold
  output_style_red
  echo "$@"
  output_style_reset
}

function echo_n_red_bold {
  output_style_red
  output_style_bold
  echo -n "$@"
  output_style_reset
}


# Prints the provided usage message
function echo_usage {
  echo_indented "$@"
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


# Pipe to this function to indent output (2 spaces by default)
function indent {
  local count=$1
  if [ -z "$1" ]; then count=2; fi
  local spaces="$(printf "%${count}s")"

  sed "s/^/${spaces}/"
}


function output_style_bold {
  tput bold
}


function output_style_cyan {
  tput setaf 6
}


function output_style_dim {
  tput dim
}


function output_style_red {
  tput setaf 1
}


function output_style_reset {
  tput sgr0
}


# Prints a command
function print_git_command {
  local branch_name=$(get_current_branch_name)
  echo_header "[$branch_name] $*"
}


function prompt_yn {
  echo -n " [Y/n] "
  read yn
  case "$yn" in
    [Yy]) return 0;;
    [Nn]) return 1;;
    *([[:space:]])) return 0;;
    *) echo "Please answer yes (y) or no (n)."; return 1;;
  esac
}


# Run a normal (non Git) command.
#
# Prints the command and the output
function run_command {
  local cmd="$1"
  echo_header "$cmd"
  eval "$cmd" 2>&1
}


# Run a Git command
#
# Prints the command and the Git branch it is running on, as well as the output.
function run_git_command {
  local cmd="$1"
  print_git_command "$cmd"
  eval "$cmd" 2>&1
}
