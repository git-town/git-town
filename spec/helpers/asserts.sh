# Provides custom-tailored asserts for Git testing


# Asserts that the given actual value is equal to the given expected value
function assert {
  actual_value=$1
  expected_value=$2
  name=$3

  if [ ! "$actual_value" = "$expected_value" ]; then
    if [ -z $name ]; then
      echo_failure "Expected '$expected_value', but got '$actual_value'"
    else
      echo_failure "Expected '$name' to equal '$expected_value', but it was '$actual_value'"
    fi
  else
    if [ -z $name ]; then
      echo_success "Found '$actual_value' as expected"
    else
      echo_success "'$name' is as expected '$actual_value'"
    fi
  fi
}


# Asserts that the branch with the given name has all of the given commits,
# in the given order.
function assert_branch_has_commit {
  if [ `git log | grep $1 | wc -l` = 1 ]; then
    echo_success "Branch has as expected commit '$1'"
  else
    echo_failure "Branch does not have commit '$1'"
  fi
}


# Asserts that a local branch with the given name exists
function assert_local_branch_exists {
  local expected_branch_name=$1
  if [ $((`git branch | grep $1 | wc -l`)) == 0 ]; then
    echo_failure "Expected local branch '$1'"
  else
    echo_success "Local branch '$1' exists"
  fi
}


# Asserts that the current branch has the given name
function assert_current_branch_is {
  determine_current_branch_name
  if [ $current_branch_name = $1 ]; then
    echo_success "The current branch is as expected '$current_branch_name'"
  else
    echo_failure "Expected the current branch to be '$1', but it was '$current_branch_name'"
  fi
}


# Asserts that there is no remote branch with the given name
function assert_no_remote_branch_exists {
  if [ `git branch -a | grep remotes/origin/$1 | wc -l` == 0 ]; then
    echo_success "There is as expected no remote branch '$1'"
  else
    echo_failure "Found a remote branch '$1'"
  fi
}


# Asserts that the workspace contains the given uncommitted changes
function assert_uncommitted_changes {
  local filename=$1
  status=`git status --porcelain | grep $filename`
  if [ "$status" == "?? $filename" ]; then
    echo_success "Found '$filename' with uncommitted changes as expected"
  else
    echo_failure "Expected '$filename' to have uncommitted changed, but it didn't"
  fi
}
