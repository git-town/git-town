# Helper methods for dealing with configuration.


# The file name of the configuration file.
#
# This file just contains the name of the main development branch.
# Typically this is either 'master' or 'development'.
config_filename=".main_branch_name"
if [ -z $config_path ]; then
  config_path=$config_filename
fi


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
  echo $main_branch_name > $config_path
  echo
  echo "I have created this file with content '$main_branch_name' for you."
  echo Please add this file to your .gitignore,
  echo then run this script again to continue.
  exit_with_error
}


# Makes sure the branch config file exists.
#
# Creates one if not with data queried from the user.
#
# Exits the script if the config file didn't exist.
function ensure_config_file_exists {
  if [[ ! -f $config_path ]]; then
    echo_error_header
    echo "  Didn't find the $config_filename file."
    echo
    create_config_file
  fi
}



ensure_config_file_exists
