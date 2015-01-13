#!/usr/bin/env bash


# Returns the url for the remote
function remote_url {
  git remote -v | grep "origin.*fetch" | awk '{print $2}'
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
