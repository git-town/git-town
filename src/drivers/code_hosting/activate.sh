#!/usr/bin/env bash

# This is the activation script for the "code hosting" driver family.
#
# It is called when this family is activated.
# It automatically determines the right driver for the current environment
# and loads it.


# Loads the code-hosting driver that works with the given hostname
function activate_driver_for_code_hosting {
  local origin_hostname="$(remote_domain)"
  if [ "$(contains_string "$origin_hostname" github)" == true ]; then
    activate_driver 'code_hosting' 'github'
  elif [ "$(contains_string "$origin_hostname" bitbucket)" == true ]; then
    activate_driver 'code_hosting' 'bitbucket'
  else
    echo_error_header
    echo_usage "Unsupported hosting service."
    echo_usage 'This command requires hosting on GitHub or Bitbucket.'
    exit_with_error newline
  fi
}
