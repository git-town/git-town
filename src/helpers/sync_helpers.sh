#!/usr/bin/env bash

# Helper functions for syncing branches


# Outputs the steps to sync the given branch
function sync_branch_steps {
  local branch=$1
  local is_feature ; is_feature="$(is_feature_branch "$branch")"

  # If there is a remote origin, then checkout and sync all branches because
  # there may be changes to perennial branches, otherwise only sync feature
  # branches because perennial branches will not need syncing
  if [ "$HAS_REMOTE" = true ] || [ "$is_feature" = true ]; then
    echo "checkout $branch"

    if [ "$is_feature" = true ]; then
      echo "merge_tracking_branch"
      echo "merge $(parent_branch "$branch")"
    else
      echo "rebase_tracking_branch"
    fi

    echo_if_true "push" "$HAS_REMOTE"
  fi
}
