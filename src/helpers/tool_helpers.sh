#!/usr/bin/env bash

# Helper methods for dealing with tools.


# Returns whether or not the tool is available
function has_tool {
  local tool=$1
  if [ "$(which "$tool" | wc -l | tr -d ' ')" = 0 ]; then
    echo false
  else
    echo true
  fi
}
