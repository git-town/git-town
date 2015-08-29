#!/usr/bin/env bash


# Pushes tags to the remote
function push_tags {
  run_git_command "git push --tags"
}
