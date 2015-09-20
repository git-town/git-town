#!/usr/bin/env bash


# Returns the url for the remote
function remote_url {
  mock_remote=$(git config --get git-town.testing.remote-url)
  if [ ! -z "$mock_remote" ]; then
    echo "$mock_remote"
  else
    git remote -v | grep "origin.*fetch" | awk '{print $2}'
  fi
}


# Returns the domain of the remote repository
function remote_domain {
  remote_url | sed -E "s#(https?://([^@]*@)?|git@)([^/:]+).*#\3#"
}


# Returns the USER/REPO for the remote repository
function remote_repository_name {
  local domain=$(remote_domain)
  remote_url | sed -E "s#.*$domain[/:](.+)#\1#" | sed "s/\.git$//"
}


# Returns true if the repository has a remote configured
function has_remote_url {
  if [ -z "$(remote_url)" ]; then
    echo false
  else
    echo true
  fi
}
