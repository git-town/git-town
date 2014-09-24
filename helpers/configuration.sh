# Helper methods for dealing with configuration.


# Stores the current name branch name in the config file.
function store_main_branch_name {
  git config git-town.main-branch-name $main_branch_name
}


# Update old configuration to new one if it exists
if [[ -f ".main_branch_name" ]]; then
  main_branch_name=`cat .main_branch_name`
  store_main_branch_name
  rm .main_branch_name
fi


# Read main branch name from config, ask and store it if it isn't known yet.
main_branch_name=`git config --get git-town.main-branch-name`
if [[ -z "$main_branch_name" ]]; then
  echo "Please enter the name of the main dev branch (typically 'master' or 'development'):"
  read main_branch_name
  if [[ -z "$main_branch_name" ]]; then
    echo_error_header
    echo "  You have not provided the name for the main branch."
    echo "  This information is necessary to run this script."
    echo "  Please try again."
    exit_with_error
  fi
  store_main_branch_name
  echo
  echo "I have stored the main branch name '$main_branch_name' for you."
  echo
  echo "I am stopping right now so that you know about this."
  echo "Please repeat the last Git Town command to execute it."
  exit
fi
