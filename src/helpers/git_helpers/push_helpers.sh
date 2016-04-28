#!/usr/bin/env bash

function hack_should_push {
  if [ "$HACK_PUSH_FLAG" == 'true' ]; then
    echo true
  else
    echo false
  fi
}


# Pushes tags to the remote
function push_tags {
  run_command "git push --tags"
}
