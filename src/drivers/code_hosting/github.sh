#!/usr/bin/env bash


function create_pull_request {
  local repository=$1
  local branch=$2
  local parent_branch=$3

  open_browser "https://github.com/$repository/compare/$parent_branch...$branch?expand=1"
}


function show_repo {
  local repository=$1

  open_browser "https://github.com/$repository"
}
