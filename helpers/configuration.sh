# Helper methods for dealing with configuration.


# The file name of the configuration file.
#
# This file just contains the name of the main development branch.
# Typically this is either 'master' or 'development'.
config_filename=".gittownrc"
old_config_filename=".main_branch_name"


# Creates the configuration file with data asked from the user.
function create_config_file {
  echo "Please enter the name of the main dev branch (typically 'master' or 'development'):"
  read main_branch_name
  if [[ -z "$main_branch_name" ]]; then
    echo_error_header
    echo "  You have not provided the name for the main branch."
    echo "  This information is necessary to run this script."
    echo "  Please try again."
    exit_with_error
  fi
  echo $main_branch_name > $config_filename
  echo
  echo "I have created this file with content '$main_branch_name' for you."
}


function update_config_file {
  echo_header "I have found an old Git Town configuration file: $old_config_filename"
  echo "I am updating it to the new format: $config_filename"
  mv $old_config_filename $config_filename
  exit_with_error
}


# Makes sure the branch config file exists.
#
# Creates one if not with data queried from the user.
#
# Exits the script if the config file didn't exist.
function ensure_config_file_exists {
  if [[ ! -f $config_filename ]]; then
    if [[ -f $old_config_filename ]]; then
      update_config_file
    else
      echo_error_header
      echo "  Didn't find the $config_filename file."
      echo
      create_config_file
    fi
  fi
}


# Stores the current name branch name in the config file.
function store_main_branch_name {
  echo $main_branch_name > $config_filename
}



ensure_config_file_exists
