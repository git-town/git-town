# Helper methods for dealing with abort/continue scripts


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
