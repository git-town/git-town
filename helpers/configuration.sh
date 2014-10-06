# Helper methods for dealing with configuration.


# Stores the current name branch name in the config file.
function store_main_branch_name {
  git config git-town.main-branch-name $main_branch_name
}


# Stores the current name branch name in the config file.
function store_special_branch_names {
  git config git-town.special-branch-names $special_branch_names
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
fi


# Read special branch names from config, ask and store if needed
special_branch_names=`git config --get git-town.special-branch-names`
if [[ $? == '1' ]]; then
  echo "Please enter the name of any special branches as a comma seperated list."
  echo "git sync and git ship will be automatically prevented on these branches."
  echo "(ex: 'qa, production')"
  read special_branch_names
  store_special_branch_names
  echo "Special branch names stored as '$special_branch_names'"
fi
