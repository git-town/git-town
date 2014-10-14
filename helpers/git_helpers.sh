# Helper methods for working with Git.


# Checks out the branch with the given name.
#
# Skips this operation if the requested branch
# is already checked out.
function checkout_branch {
  local branch_name=$1
  if [ ! "`get_current_branch_name`" = "$branch_name" ]; then
    run_command "git checkout $branch_name"
  fi
}


# Checks out the main development branch in Git.
#
# Skips the operation if we already are on that branch.
function checkout_main_branch {
  checkout_branch $main_branch_name
}


# Cherry picks the shas into the current branch
function cherry_pick {
  local shas=$*
  run_command "git cherry-pick $shas"
  if [ $command_exit_status != 0 ]; then error_cherry_pick $shas; fi
}

# Creates a new feature branch with the given name.
#
# The feature branch is cut off the main development branch.
function create_feature_branch {
  local new_branch_name=$1
  run_command "git checkout -b $new_branch_name $main_branch_name"
}


# Deletes the given branch from both the local machine and on remote.
function delete_branch {
  local branch_name=$1
  local current_branch_name=`get_current_branch_name`
  checkout_branch $branch_name
  if [ `has_tracking_branch` == true ]; then
    run_command "git push origin :${branch_name}"
  fi
  checkout_branch $current_branch_name
  run_command "git branch -D $branch_name"
}


# Determines whether there are open changes in Git.
function has_open_changes {
  if [ `git status --porcelain | wc -l` == 0 ]; then
    echo false
  else
    echo true
  fi
}


# Determines whether the current branch has a remote tracking branch.
function has_tracking_branch {
  local current_branch_name=`get_current_branch_name`
  if [ `git branch -vv | grep "\* $current_branch_name\b" | grep "\[origin\/$current_branch_name.*\]" | wc -l` == 0 ]; then
    echo false
  else
    echo true
  fi
}


# Exists the application with an error message if the
# current working directory contains uncommitted changes.
function ensure_no_open_changes {
  if [ $initial_open_changes = true ]; then
    echo_error_header
    echo $*
    exit_with_error
  fi
}


# Exists the application with an error message if the working directory
# is on the main development branch.
function ensure_on_feature_branch {
  local error_message=$1
  if [ "`get_current_branch_name`" = "$main_branch_name" ]; then
    echo_error_header
    echo "  $error_message"
    exit_with_error
  fi
}


# Called by pull_branch when the merge/rebase fails with conflicts
function error_pull_branch {
  if [ $1 == $main_branch_name ]; then
    error_pull_main_branch
  else
    error_pull_feature_branch
  fi
}


# Returns the current branch name
function get_current_branch_name {
  git branch | grep "*" | awk '{print $2}'
}



# Fetches updates from the central repository.
#
# It is safe to call this method multiple times per session,
# since it makes sure that it fetches updates only once per session
# by tracking this through the global variable $repo_fetched.
function fetch_repo {
  if [ $repo_fetched == false ]; then
    run_command "git fetch --prune"
    repo_fetched=true
  fi
}
repo_fetched=false


# Fetches changes from the upstream repository
function fetch_upstream {
  run_command "git fetch upstream"
}


# Returns true if the repository has a branch with the given name
function has_branch {
  local branch_name=$1
  if [ `git branch | grep "$branch_name" | wc -l` == 0 ]; then
    echo false
  else
    echo true
  fi
}


# Merges the given branch into the current branch
function merge_branch {
  local branch_name=$1
  local current_branch_name=`get_current_branch_name`
  run_command "git merge $branch_name"
  if [ $command_exit_status != 0 ]; then error_merge_branch; fi
}


# Pulls updates of the feature branch from the remote repo
function pull_branch {
  local strategy=$1
  local current_branch_name=`get_current_branch_name`
  if [ -z $strategy ]; then strategy='merge'; fi
  if [ `has_tracking_branch` == true ]; then
    fetch_repo
    run_command "git $strategy origin/$current_branch_name"
    if [ $command_exit_status != 0 ]; then error_pull_branch $current_branch_name; fi
  else
    echo "Branch '$current_branch_name' has no remote branch, skipping pull of updates"
  fi
}


# Pulls updates of the current branch from the upstream repo
function pull_upstream_branch {
  local current_branch_name=`get_current_branch_name`
  fetch_upstream
  run_command "git rebase upstream/$current_branch_name"
  if [ $command_exit_status != 0 ]; then error_pull_upstream_branch; fi
}


# Pushes the branch with the given name to origin
function push_branch {
  local current_branch_name=`get_current_branch_name`
  if [ `has_tracking_branch` = true ]; then
    run_command "git push"
  else
    run_command "git push -u origin $current_branch_name"
  fi
}


# Pushes tags to the remote
function push_tags {
  run_command "git push --tags"
}


# Returns the url for the remote with the specified name
function remote_url {
  git remote -v | grep "$1.*fetch" | awk '{print $2}'
}


# Unstashes changes that were stashed in the beginning of a script.
#
# Only does this if there were open changes when the script was started.
function restore_open_changes {
  if [ $initial_open_changes = true ]; then
    run_command "git stash pop"
  fi
}


# Squash merges the given branch into the current branch
function squash_merge {
  local branch_name=$1
  local commit_message=$2
  local current_branch_name=`get_current_branch_name`
  run_command "git merge --squash $branch_name"
  if [ $command_exit_status != 0 ]; then error_squash_merge; fi
  if [ "$commit_message" == "" ]; then
    print_command "git commit -a"
    git commit -a
  else
    print_command "git commit -a -m \"$commit_message\""
    git commit -a -m "$commit_message"
  fi
}


# Stashes uncommitted changes if they exist.
function stash_open_changes {
  if [ $initial_open_changes = true ]; then
    run_command "git stash -u"
  fi
}

# Stashes uncommitted changes if they exist.
function sync_main_branch {
  local current_branch_name=`get_current_branch_name`
  checkout_main_branch
  pull_branch 'rebase'
  checkout_branch $current_branch_name
}
