#!/usr/bin/env bash


function revert_commit {
  local commit=$1
  run_command "git revert $commit"
}
