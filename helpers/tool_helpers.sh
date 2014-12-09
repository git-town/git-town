#!/bin/bash

# Helper methods for dealing with tools.


# Ensures that the given tool is installed.
function ensure_tool_installed {
  local toolname=$1
  if [ "$(which "$toolname" | wc -l)" == 0 ]; then
    echo_error_header
    echo_error "You need the '$toolname' tool in order to run tests."
    echo_error "Please install it using your package manager,"
    echo_error "or on OS X with 'brew install $toolname'."
    exit_with_error
  fi
}
