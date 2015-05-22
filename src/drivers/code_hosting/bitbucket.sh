#!/usr/bin/env bash


# Returns the source for a new Bitbucket pull request
function bitbucket_source {
  local repository=$1
  local branch=$2

  local sha=$(git log --format="%H" -1 | cut -c-12)
  echo "$repository:$sha:$branch" | sed 's/\//%2F/g' | sed 's/\:/%3A/g'
}


function create_pull_request {
  local repository=$1
  local branch=$2

  open_browser "https://bitbucket.org/$repository/pull-request/new?source=$(bitbucket_source "$repository" "$branch")"
}


function show_repo {
  local repository=$1

  open_browser "https://bitbucket.org/$repository"
}
