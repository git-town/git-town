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
      echo "${PULL_BRANCH_STRATEGY}_tracking_branch"

      if [ "$branch" = "$MAIN_BRANCH_NAME" ] && [ "$(has_remote_upstream)" = true ]; then
        echo "fetch_upstream"
        echo "rebase upstream/$MAIN_BRANCH_NAME"
      fi
    fi

    if [ "$HAS_REMOTE" = true ]; then
      if [ "$(has_tracking_branch "$branch")" == true ]; then
        echo "push_branch $branch"
      else
        echo "create_tracking_branch $branch"
      fi
    fi
  fi
}
