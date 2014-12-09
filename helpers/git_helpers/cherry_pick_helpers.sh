#!/bin/bash


# Abort a cherry-pick
function abort_cherry_pick {
  run_command "git cherry-pick --abort"
}


# Cherry picks the SHAs into the current branch
function cherry_pick {
  local SHAs="$*"
  run_command "git cherry-pick $SHAs"
  if [ $? != 0 ]; then error_cherry_pick "$SHAs"; fi
}
