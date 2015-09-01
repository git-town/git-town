#!/usr/bin/env bash


# Pushes tags to the remote
function push_tags {
  run_command "git push --tags"
}
