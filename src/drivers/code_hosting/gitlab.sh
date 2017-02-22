#!/usr/bin/env bash


function create_pull_request {
  local repository=$1
  local src="merge_request%5Bsource_branch%5D=${2}"
  local dest="merge_request%5Btarget_branch%5D=${3:-master}"

  open_browser "https://gitlab.com/$repository/merge_requests/new?${src}\&${dest}"
}


function show_repo {
  local repository=$1

  open_browser "https://gitlab.com/$repository"
}
