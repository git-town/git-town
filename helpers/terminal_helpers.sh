#!/bin/bash

# Helper methods for writing to the terminal.


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
  local str=`echo_indented Error`
  echo
  echo_red_bold "$str"
}


# Prints the provided error message
function echo_error {
  local str=`echo_indented "$*"`
  echo_red "$str"
}


# Prints the message idented
function echo_indented {
  echo "  $*"
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
  local str=`echo_indented Usage`
  echo_header "$str"
}


# Exits the currently running script with an error response code.
function exit_with_error {
  echo
  exit 1
}


# Exits the currently running script with a success response code
function exit_with_success {
  echo
  exit 0
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
  local branch_name=`get_current_branch_name`
  echo_header "[$branch_name] $*"
}


# Run a command, prints command and output
function run_command {
  local cmd=$*
  print_command $cmd
  eval $cmd 2>&1
}
