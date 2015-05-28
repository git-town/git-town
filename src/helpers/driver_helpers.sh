#!/usr/bin/env bash

# Helper methods around drivers


# Loads and activates the given driver
function activate_driver {
  local driver_family_name=$1
  local driver_name=$2

  source "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/../drivers/$driver_family_name/$driver_name.sh"
}


# Activates the given driver family,
# and lets the driver family activation script determine on its own
# which particular driver to load.
function activate_driver_family {
  local driver_family_name=$1

  source "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/../drivers/$driver_family_name/activate.sh"
  shift
  activate_driver_for_"$driver_family_name" "$@"
}
