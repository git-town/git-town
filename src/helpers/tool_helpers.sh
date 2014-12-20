#!/bin/bash

# Helper methods for dealing with tools.

# Ensures that the given tool is installed.
function ensure_tool_installed {
  local toolname=$1
  if [ "$(has_tool "$toolname")" = false ]; then
    echo_error_header
    echo_error "You need the '$toolname' tool in order to run tests."
    echo_error "Please install it using your package manager,"
    echo_error "or on OS X with 'brew install $toolname'."
    exit_with_error
  fi
}


# Returns whether or not the tool is available
function has_tool {
  local tool=$1
  if [ "$(which "$tool" | wc -l | tr -d ' ')" = 0 ]; then
    echo false
  else
    echo true
  fi
}
