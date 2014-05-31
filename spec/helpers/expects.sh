# Provides custom-tailored asserts for Git testing


# Asserts that there is currently a cherry-picking process going on in the current branch
function expect_cherrypick_in_progress {
  if [ `git status | grep "You are currently cherry-picking" | wc -l` == 1 ]; then
    echo_success "A cherry-pick is currently active"
  else
    echo_failure "There is no cherry-pick active right now"
  fi
}


# Asserts that there is currently a cherry-picking process going on in the current branch
function expect_no_cherrypick_in_progress {
  if [ `git status | grep "You are currently cherry-picking" | wc -l` == 0 ]; then
    echo_success "There is no cherry-pick active right now"
  else
    echo_failure "A cherry-pick is currently active"
  fi
}


# Asserts that the given file has merge conflicts
function expect_conflict_for_file {
  if [ `git status | grep "both added.*$1" | wc -l` == 1 ]; then
    echo_success "File '$1' has as expected merge conflicts"
  else
    echo_failure "File '$1' does not have merge conflicts"
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


# Expects that the file with the given path exists
function expect_file_exists {
  file_name=$1
  if [ `ls $file_name | wc -l` == 1 ]; then
    echo_success "File '$file_name' exists"
  else
    echo_failure "File '$file_name' does not exist"
  fi
}


# Expects that the file with the given path does not exist
function expect_file_does_not_exist {
  file_name=$1
  if [ `ls $file_name | wc -l` == 0 ]; then
    echo_success "File '$file_name' does not exist"
  else
    echo_failure "File '$file_name' does exist"
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


# Asserts that the last commit includes the given text
function expect_last_commit_includes {
  if [ `git log | grep $1 | wc -l` == 0 ]; then
    echo_failure "Expected the last commit to contain '$1', but it didn't"
  else
    echo_success "The last commit contains as expected '$1'"
  fi
}


# Asserts that there are exactly the given number of local branches
function expect_local_branch_count {
  branch_count=`git branch | wc -l`
  if [ $branch_count == $1 ]; then
    echo_success "As expected there are $1 local branches"
  else
    echo_failure "Expected $1 local branches, but got $branch_count"
  fi
}


# Asserts that the branch with the given name has all of the given commits,
# in the given order.
function expect_local_branch_has_commit {
  branch_name=$1
  commit_name=$2
  checkout_branch $branch_name
  if [ `git log | grep $commit_name | wc -l` = 0 ]; then
    echo_failure "Branch '$branch_name' does not have commit '$commit_name'"
  else
    echo_success "Branch '$branch_name' has as expected commit '$commit_name'"
  fi
}


# Asserts that a local branch with the given name exists
function expect_no_local_branch_exists {
  local expected_branch_name=$1
  if [ `git branch | grep $1 | wc -l` == 0 ]; then
    echo_success "Local branch '$1' does not exist"
  else
    echo_failure "Local branch '$1' does exist"
  fi
}


# Asserts that a local branch with the given name exists
function expect_local_branch_exists {
  local expected_branch_name=$1
  if [ `git branch | grep $1 | wc -l` == 0 ]; then
    echo_failure "Expected local branch '$1'"
  else
    echo_success "Local branch '$1' exists"
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


# Asserts that there is currently a rebase in progress
function expect_no_rebase_in_progress {
  determine_rebase_in_progress
  if [ $rebase_in_progress == false ]; then
    echo_success "Currently no rebase in progress"
  else
    echo_failure "Currently a rebase in progress"
  fi
}


# Asserts that there is currently a rebase in progress
function expect_rebase_in_progress {
  determine_rebase_in_progress
  if [ $rebase_in_progress == true ]; then
    echo_success "Currently a rebase in progress"
  else
    echo_failure "Currently no rebase in progress"
  fi
}


# Asserts that the given branch is fully synchronized with its remote branch
function expect_synchronized_branch {
  if [ `git branch -vv | grep $1 | grep -o '\[.*\]' | tr -d '[]' | awk '{ print $2 }' | tr -d '\n' | wc -m` == 0 ]; then
    echo_success "Branch '$1' is fully synchronized with its remote branch"
  else
    echo_failure "Branch '$1' is not completely synchronized with its remote branch"
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

