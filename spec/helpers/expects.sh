# Provides custom-tailored asserts for Git testing


# Asserts that the branch with the given name has all of the given commits,
# in the given order.
function expect_branch_has_commit {
  if [ `git log | grep $1 | wc -l` = 1 ]; then
    echo_success "Branch has as expected commit '$1'"
  else
    echo_failure "Branch does not have commit '$1'"
  fi
}


# Asserts that a local branch with the given name exists
function expect_local_branch_exists {
  local expected_branch_name=$1
  if [ $((`git branch | grep $1 | wc -l`)) == 0 ]; then
    echo_failure "Expected local branch '$1'"
  else
    echo_success "Local branch '$1' exists"
  fi
}


# Asserts that the current branch has the given name
function expect_current_branch_is {
  determine_current_branch_name
  if [ $current_branch_name = $1 ]; then
    echo_success "The current branch is as expected '$current_branch_name'"
  else
    echo_failure "Expected the current branch to be '$1', but it was '$current_branch_name'"
  fi
}


# Expects that the file with the given name has the given content
function expect_file_content {
  file_name=$1
  expected_file_content=$2
  actual_file_content=`cat $file_name`
  if [ "$actual_file_content" == "$expected_file_content" ]; then
    echo_success "File '$file_name' has as expected content '$expected_file_content'"
  else
    echo_failure "File '$file_name' was expected to have content '$expected_file_content', but has '$actual_file_content'"
  fi
}

# Asserts that there is no remote branch with the given name
function expect_no_remote_branch_exists {
  if [ `git branch -a | grep remotes/origin/$1 | wc -l` == 0 ]; then
    echo_success "There is as expected no remote branch '$1'"
  else
    echo_failure "Found a remote branch '$1'"
  fi
}


# Asserts that the workspace contains the given uncommitted changes
function expect_uncommitted_changes {
  local filename=$1
  status=`git status --porcelain | grep $filename`
  if [ "$status" == "?? $filename" ]; then
    echo_success "Found '$filename' with uncommitted changes as expected"
  else
    echo_failure "Expected '$filename' to have uncommitted changed, but it didn't"
  fi
}
