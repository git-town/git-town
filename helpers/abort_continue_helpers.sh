# Helper methods for dealing with abort/continue scripts


function create_pull_main_branch_abort_script {
  echo "git rebase --abort" > $abort_script_filename
  echo "git checkout $initial_branch_name" >> $abort_script_filename
  if [ $initial_open_changes == true ]; then
    echo "git stash pop" >> $abort_script_filename
  fi
}


function create_pull_feature_branch_abort_script {
  echo "git merge --abort" > $abort_script_filename
  if [ $initial_open_changes == true ]; then
    echo "git stash pop" >> $abort_script_filename
  fi
}


function create_merge_main_branch_abort_script {
  echo "git merge --abort" > $abort_script_filename
  if [ $initial_open_changes == true ]; then
    echo "git stash pop" >> $abort_script_filename
  fi
}


function exit_with_abort_continue_messages {
  local cmd=$1

  echo
  if [ -n $abort_script_filename ]; then
    echo_red "To abort, run \"git $cmd --abort\"."
  fi
  if [ -n $continue_script_filename ]; then
    echo_red "To continue after you have resolved the conflicts, run \"git $cmd --continue\"."
  fi
  exit_with_error
}


function remove_abort_continue_script_filenames {
  if [ -n $abort_script_filename ]; then rm $abort_script_filename; fi
  if [ -n $continue_script_filename ]; then rm $continue_script_filename; fi
}


function run_abort_script {
  if [ -f $abort_script_filename ]; then
    source $abort_script_filename
    remove_abort_continue_script_filenames
  else
    echo_red "Cannot find abort definition file"
  fi
}


function run_continue_script {
  if [ -f $continue_script_filename ]; then
    source $continue_script_filename
    remove_abort_continue_script_filenames
  else
    echo_red "Cannot find continue definition file"
  fi
}
