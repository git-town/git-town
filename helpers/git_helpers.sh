# Helper methods for working with Git.


# Checks out the branch with the given name.
#
# Skips this operation if the requested branch
# is already checked out.
function checkout_branch {
  local branch_name=$1
  if [ ! "`get_current_branch_name`" = "$branch_name" ]; then
    git checkout $branch_name
  fi
}


# Checks out the main development branch in Git.
#
# Skips the operation if we already are on that branch.
function checkout_main_branch {
  checkout_branch $main_branch_name
}


# Creates a new feature branch with the given name.
#
# The feature branch is cut off the main development branch.
function create_feature_branch {
  local new_branch_name=$1
  echo_header "Creating the '$new_branch_name' branch off the '$main_branch_name' branch"
  git checkout -b $new_branch_name $main_branch_name
}


# Deletes the branch with the given name
# Will be on the main branch when the function ends
function delete_branch {
  local branch_name=$1
  echo_header "Removing the old '$branch_name' branch"
  checkout_branch $branch_name
  local tracking_branch_exists=`has_tracking_branch`
  checkout_main_branch
  git branch -D $branch_name
  if [ $tracking_branch_exists = true ]; then
    git push origin :${branch_name}
  fi
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
    echo_header "  Error"
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
    git fetch --prune
    repo_fetched=true
  fi
}
repo_fetched=false


# Fetches changes from the upstream repository
function fetch_upstream {
  echo_header "Fetching updates for 'upstream'"
  git fetch upstream
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


# Merges the main branch into the current branch
function merge_main_branch {
  local branch_name=`get_current_branch_name`
  echo_header "Merging '$main_branch_name' into '$branch_name'"
  git merge $main_branch_name
  if [ $? != 0 ]; then error_merge_main_branch $branch_name; fi
}


# Pulls updates of the current branch from origin
function pull_branch {
  local branch_name=`get_current_branch_name`
  echo_header "Pulling updates for the '$branch_name' branch"
  if [ `has_tracking_branch` == true ]; then
    fetch_repo
    git merge origin/$branch_name
    if [ $? != 0 ]; then error_pull_branch $branch_name; fi
  else
    echo "Branch '$branch_name' has no remote branch, skipping pull of updates"
  fi
}


# Pulls updates for the main branch
# Will be on the same branch as when the function begins
function pull_main_branch {
  local branch_name=`get_current_branch_name`
  checkout_main_branch
  pull_branch
  checkout_branch $branch_name
}


# Pulls updates of the current branch from the upstream repo
function pull_upstream {
  local branch_name=`get_current_branch_name`
  echo_header "Pulling 'upstream/$branch_name' into '$branch_name'"
  fetch_upstream
  git merge upstream/$branch_name
}


# Pushes the current branch to origin
function push_branch {
  local branch_name=`get_current_branch_name`
  echo_header "Pushing '$branch_name'"
  if [ `has_tracking_branch` = true ]; then
    git push
  else
    git push -u origin $branch_name
  fi
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
    echo_header "Restoring uncommitted changes"
    git stash pop
  fi
}


# Merges the current branch into the main branch
# Ends on the main branch
function squash_merge_into_main_branch {
  local commit_message=$1
  local branch_name=`get_current_branch_name`
  echo_header "Merging the '$branch_name' branch into '$main_branch_name'"
  checkout_main_branch
  if [ "$commit_message" == "" ]; then
    git merge --squash $branch_name && git commit -a
  else
    git merge --squash $branch_name && git commit -a -m "$commit_message"
  fi
  if [ $? != 0 ]; then
    error_squash_merge_into_main_branch $branch_name;
  fi
}


# Stashes uncommitted changes if they exist.
function stash_open_changes {
  if [ $initial_open_changes = true ]; then
    echo_header "Stashing uncommitted changes"
    git add -A
    git stash
  fi
}



