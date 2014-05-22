# Tools for managing the Git repo used in the specs.


current_time=`date +%s`

# Name of the main branch used in the current spec runner session
main_branch_name="main_$$_$current_time"

# Name of the primary feature branch used in the current spec runner session
feature_branch_name="feature_$$_$current_time"

# Path to the directory that contains the test repo
test_repo_parent_dir="/tmp"

# Name of the Git repo used to run the tests
test_repo_name="git_town_specs"

# Full path to the Git repo used to run the tests
test_repo_path="$test_repo_parent_dir/$test_repo_name"



# Adds a new commit with the given changes to the given local branch
function add_local_commit {
  local branch_name=$1
  local commit_msg=$2
  local file_name=$3
  if [ -z $file_name ]; then file_name=$commit_msg; fi
  local content=$4
  if [ -z $content ]; then content=$commit_msg; fi
  echo "adding local commit $commit_msg to branch $branch_name with file $file_name and content $content"

  checkout_branch $branch_name
  echo $content >> $file_name
  git add $file_name
  git commit -m $commit_msg
}


# Adds a new commit to the remote repo only, not to the local repo
function add_remote_commit {
  local branch_name=$1
  local commit_msg=$2
  local file_name=$3
  if [ -z $file_name ]; then file_name=$commit_msg; fi
  local content=$4
  if [ -z $content ]; then content=$commit_msg; fi
  echo "adding local commit $commit_msg to branch $branch_name with file $file_name and content $content"

  checkout_branch $branch_name
  echo $content >> $file_name
  git add $file_name
  git commit -m $commit_msg
  git push
  git reset --hard HEAD^
}


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
    git clone git@github.com:Originate/git_town_specs.git
  fi

  # Enter repo directory
  cd $test_repo_name
}



# Removes all branches from the remote repo
function remove_all_remote_branches {
  echo "Removing all remote branches"
  git branch -a            | # get all branches
    grep remotes/origin    | # keep only the remote ones
    grep -v HEAD           | # remove HEAD branch
    grep -v master         | # remove master branch
    cut -d '/' -f 3        | # keep only the third part of the path
    awk '{ print ":"$0 }'  | # append ':' before each branch name
    xargs git push origin    # delete the branches
}


# Removes all branches from the remote repo that currently exist on the local machine
function remove_all_my_remote_branches {
  echo "Removing all my remote branches"
  checkout_branch "master"
  git branch -vv          |
    grep -o '\[.*\]'      |
    grep -v master        |
    tr -d '[]'            |
    cut -d '/' -f 2       |
    awk '{ print ":"$0 }' |
    xargs git push origin
}


# Removes all local branches except master
function remove_all_local_branches {
  echo "Removing all local branches"
  git branch | grep -v master | tr -d '*' | xargs git branch -D
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
  if [ $remote_main_branch_created = false ]; then
    echo "creating a remote main branch"
    git push -u origin $main_branch_name
    remote_main_branch_created=true
  fi
}


# Removes all branches from the test repo
function reset_test_repo {
  checkout_main_branch
  remove_all_my_remote_branches
  remove_all_local_branches
}

