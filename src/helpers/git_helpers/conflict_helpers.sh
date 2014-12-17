#!/bin/bash


# Exits if there are unresolved conflicts
function ensure_no_conflicts {
  if [ "$(has_conflicts)" == true ]; then
    echo_error_header
    echo_error "You must resolve the conflicts before continuing the $git_command"
    exit_with_error
  fi
}


# Returns true if there are conflicts
function has_conflicts {
  if [ "$(git status | grep -c 'Unmerged paths')" == 0 ]; then
    echo false
  else
    echo true
  fi
}
