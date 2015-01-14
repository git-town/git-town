#!/usr/bin/env bash


# Returns the author with the most commits in the given branch that are not in the main branch
function branch_author {
  local branch_name=$1
  git log "$main_branch_name..$branch_name" --format='%an <%ae>' | # Authors of commits only in $branch_name
    uniq -c  | # get count for each author
    sort -rn | # reverse numeric sort
    head -1  | # take only the first
    sed -E 's/^ *[0-9]+ //g' # remove the count
}


# Returns the default author for new commits in the repository
function local_author {
  local name="$(git config user.name)"
  local email="$(git config user.email)"
  echo "$name <$email>"
}
