#!/usr/bin/env bash


# Returns the source for a new Bitbucket pull request
function bitbucket_source {
  local repository=$1
  local branch=$2

  local sha ; sha=$(git log --format="%H" -1 | cut -c-12)
  echo "$repository:$sha:$branch" | sed 's/\//%2F/g' | sed 's/\:/%3A/g'
}


# Returns the parent for a new Bitbucket pull request
function bitbucket_destination {
  local repository=$1
  local branch=$2

  echo "$repository::$branch" | sed 's/\//%2F/g' | sed 's/\:/%3A/g'
}


function create_pull_request {
  local repository=$1
  local branch=$2
  local parent_branch=$3

  src=$(bitbucket_source "$repository" "$branch")
  dest=$(bitbucket_destination "$repository" "$parent_branch")
  open_browser "https://bitbucket.org/$repository/pull-request/new?source=$src\&dest=$dest"
}


function show_repo {
  local repository=$1

  open_browser "https://bitbucket.org/$repository"
}
