# Tools for managing the Git repo used in the specs.


current_time=`date +%s`

# Name of the main branch used in the current spec runner session
main_branch_name="main_$$_$current_time"

# Name of the primary feature branch used in the current spec runner session
feature_branch_name="feature_$$_$current_time"

# Path to the directory that contains the test repo
test_repo_parent_dir="/tmp/git_town_specs"

# Name of the Git repo used to run the tests
test_repo_name="git_town_specs"

# Full path to the Git repo used to run the tests
test_repo_path="$test_repo_parent_dir/$test_repo_name"


# Prepares the test repo for use
function enter_test_repo {

  # Create test directory if it doesn't exist yet
  if [ ! -d "$test_repo_parent_dir" ]; then
    mkdir $test_repo_parent_dir
  fi

  # Enter test directory
  cd $test_repo_parent_dir

  # Clone repo if it doesn't exist yet
  if [ ! -d "$test_repo_path" ]; then
    echo_header "PREPARING THE TEST REPO"
    hub clone git@github.com:Originate/$test_repo_name.git
  fi

  # Enter repo directory
  cd $test_repo_name

  reset_test_repo
}



# Removes all branches from the remote repo
function remove_all_my_remote_branches {
  echo "TODO: implement"
}


# Removes all local branches except master
function remove_all_local_branches {
  echo "TODO: implement"
}


# Specs must call this method if they require a main branch
# in the repo.
function require_main_branch {
  if [ $main_branch_created = false ]; then
    git checkout -b $main_branch_name
    main_branch_created=true
  fi
}


# Specs must call this method if they require the main branch
# on Github.
function require_remote_main_branch {
  require_main_branch
  if [ $main_branch_pushed = false ]; then
    git push -u origin $main_branch_name
    main_branch_pushed=true
  fi
}


# Removes all branches from the test repo
function reset_test_repo {
  remove_all_my_remote_branches
  remove_all_local_branches
}

