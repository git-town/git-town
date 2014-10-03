# Helper methods for working with Git.


# Checks out the branch with the given name.
#
# Skips this operation if the requested branch
# is already checked out.
function checkout_branch {
  local branch_name=$1
  local current_branch_name=`get_current_branch_name`
  if [ ! "$current_branch_name" = "$branch_name" ]; then
    git checkout $branch_name
  fi
}


# Checks out the main development branch in Git.
#
# Skips the operation if we already are on that branch.
function checkout_main_branch {
  checkout_branch $main_branch_name
}


# Checks out the current feature branch in Git.
#
# Skips the operation if we already are on that branch.
function checkout_feature_branch {
  checkout_branch $feature_branch_name
}


# Creates a new feature branch with the given name.
#
# The feature branch is cut off the main development branch.
function create_feature_branch {
  local new_branch_name=$1
  echo_header "Creating the '$new_branch_name' branch off the '$main_branch_name' branch"
  git checkout -b $new_branch_name $main_branch_name
}


# Deletes the given feature branch from both the local machine and on remote.
#
# If you provide 'force' as an argument, it deletes the branch even if it has
# unmerged changes.
function delete_feature_branch {
  local options=$1
  echo_header "Removing the old '$feature_branch_name' branch"
  checkout_feature_branch
  local has_tracking_branch=`determine_tracking_branch`
  checkout_main_branch
  if [[ "$options" == "force" ]]; then
    git branch -D $feature_branch_name
  else
    git branch -d $feature_branch_name
  fi
  if [ $has_tracking_branch = true ]; then
    git push origin :${feature_branch_name}
  fi
}


# Determines whether there are open changes in Git.
function determine_open_changes {
  if [ `git status --porcelain | wc -l` == 0 ]; then
    echo false
  else
    echo true
  fi
}


# Determines whether there is currently a rebase in progress.
function determine_rebase_in_progress {
  if [ `git status | grep 'You are currently rebasing' | wc -l` == 1 ]; then
    echo true
  else
    echo false
  fi
}

# Determines whether the feature branch has a remote tracking branch.
function determine_tracking_branch {
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
  if [ $has_open_changes = true ]; then
    echo_header "  Error"
    echo $*
    exit_with_error
  fi
}


# Exists the application with an error message if the working directory
# is on the main development branch.
function ensure_on_feature_branch {
  local error_message=$1
  local current_branch_name=`get_current_branch_name`
  if [ "$current_branch_name" = "$main_branch_name" ]; then
    echo_error_header
    echo "  $error_message"
    exit_with_error
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
    git fetch -p
    repo_fetched=true
  fi
}
repo_fetched=false


# Fetches changes from the upstream repository
function fetch_upstream {
  echo_header "Fetching updates for 'upstream'"
  git fetch upstream
}


# Pulls updates of the feature branch from the remote repo
function pull_feature_branch {
  echo_header "Pulling updates for the '$feature_branch_name' branch"
  checkout_feature_branch
  local has_tracking_branch=`determine_tracking_branch`
  if [ $has_tracking_branch == true ]; then
    fetch_repo
    git rebase origin/$feature_branch_name
    if [ $? != 0 ]; then error_pull_feature_branch; fi
  else
    echo "Branch '$feature_branch_name' has no remote branch, skipping pull of updates"
  fi
}


# Updates the current development branch.
function pull_main_branch {
  echo_header "Pulling updates for the '$main_branch_name' branch"
  checkout_main_branch
  local has_tracking_branch=`determine_tracking_branch`
  if [ $has_tracking_branch == true ]; then
    fetch_repo
    git rebase origin/$main_branch_name
  else
    echo "Branch '$main_branch_name' has no remote"
  fi
}


# Pulls updates of the current branch fromt the upstream repo
function pull_upstream {
  local current_branch_name=`get_current_branch_name`
  echo_header "Pulling 'upstream/$current_branch_name' into '$current_branch_name'"
  fetch_upstream
  git merge upstream/$current_branch_name
}


# Pushes the branch with the given name to origin
function push_branch {
  local branch_name=$1
  checkout_branch $branch_name
  echo_header "Pushing '$branch_name'"
  local has_tracking_branch=`determine_tracking_branch`
  if [ $has_tracking_branch = true ]; then
    git push
  else
    git push -u origin $branch_name
  fi
}


# Pushes the current feature branch to origin
function push_feature_branch {
  local options=$1
  echo_header "Pushing the updated '$feature_branch_name' to the repo"
  local has_tracking_branch=`determine_tracking_branch`
  if [ $has_tracking_branch == true ]; then
    git push
  else
    git push -u origin $feature_branch_name
  fi
}


# Pushes the main branch to origin
function push_main_branch {
  push_branch $main_branch_name
}


# Returns the url for the remote with the specified name
function remote_url {
  local url=`git remote -v | grep "$1.*fetch" | awk '{print $2}'`
  echo $url
}


# Unstashes changes that were stashed in the beginning of a script.
#
# Only does this if there were open changes when the script was started.
function restore_open_changes {
  if [ $has_open_changes = true ]; then
    echo_header "Restoring uncommitted changes"
    git stash pop
  fi
}


# Merges the current feature branch into the main dev branch.
function squash_merge_feature_branch {
  local commit_message=$1
  echo_header "Merging the '$feature_branch_name' branch into '$main_branch_name'"
  checkout_main_branch
  if [ "$commit_message" == "" ]; then
    git merge --squash $feature_branch_name && git commit -a
  else
    git merge --squash $feature_branch_name && git commit -a -m $*
  fi
  if [ $? != 0 ]; then error_squash_merge_feature_branch; fi
}


# Stashes uncommitted changes if they exist.
function stash_open_changes {
  if [ $has_open_changes = true ]; then
    echo_header "Stashing uncommitted changes"
    git add -A
    git stash
  fi
}


# Updates the current feature branch.
function update_feature_branch {
  echo_header "Rebasing the '$feature_branch_name' branch against '$main_branch_name'"
  checkout_feature_branch
  git merge $main_branch_name
  if [ $? != 0 ]; then error_update_feature_branch; fi
}
