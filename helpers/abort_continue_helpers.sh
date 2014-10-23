# Helper methods for dealing with abort/continue scripts


function create_pull_main_branch_abort_script {
  echo "git rebase --abort" > $abort_script
  echo "git checkout $initial_branch_name" >> $abort_script
  if [ $initial_open_changes == true ]; then
    echo "git stash pop" >> $abort_script
  fi
}


function create_pull_feature_branch_abort_script {
  echo "git merge --abort" > $abort_script
  if [ $initial_open_changes == true ]; then
    echo "git stash pop" >> $abort_script
  fi
}


function create_merge_main_branch_abort_script {
  echo "git merge --abort" > $abort_script
  if [ $initial_open_changes == true ]; then
    echo "git stash pop" >> $abort_script
  fi
}


function exit_with_abort_continue_messages {
  local cmd=$1

  echo
  if [ -n $abort_script ]; then
    echo_red "To abort, run \"git $cmd --abort\"."
  fi
  if [ -n $continue_script ]; then
    echo_red "To continue after you have resolved the conflicts, run \"git $cmd --continue\"."
  fi
  exit_with_error
}


function remove_abort_continue_scripts {
  if [ -n $abort_script ]; then rm $abort_script; fi
  if [ -n $continue_script ]; then rm $continue_script; fi
}


function run_abort_script {
  if [ -f $abort_script ]; then
    source $abort_script
    remove_abort_continue_scripts
  else
    echo_red "Cannot find abort definition file"
  fi
}


function run_continue_script {
  if [ -f $continue_script ]; then
    source $continue_script
    remove_abort_continue_scripts
  else
    echo_red "Cannot find continue definition file"
  fi
}
