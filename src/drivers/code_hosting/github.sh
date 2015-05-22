#!/usr/bin/env bash


function create_pull_request {
  local repository=$1
  local branch=$2

  open_browser "https://github.com/$repository/compare/$branch?expand=1"
}


function show_repo {
  local repository=$1

  open_browser "https://github.com/$repository"
}
