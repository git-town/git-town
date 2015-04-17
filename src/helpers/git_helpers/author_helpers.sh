#!/usr/bin/env bash


# Returns the author with the most commits in the given branch that are not in the main branch
function branch_author {
  local branch_name=$1
  git log "$MAIN_BRANCH_NAME..$branch_name" --format='%an <%ae>' | # Authors of commits only in $branch_name
    uniq -c  | # get count for each author
    sort -rn | # reverse numeric sort
    head -1  | # take only the first
    sed -E 's/^ *[0-9]+ //g' # remove the count
}


# Returns whether or not the given author ("user <email>") is the current git user
function is_current_user {
  local author=$1
  local author_user_name="$(echo "$author" | cut -d ' ' -f 1)"
  local local_user_name="$(git config user.name)"

  if [ "$local_user_name" = "$author_user_name" ]; then
    echo true
  else
    echo false
  fi
}
