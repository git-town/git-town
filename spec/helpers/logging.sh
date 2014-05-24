# Collects error messages and information about which tests pass and which fail.

# File path of the log file for errors
failures_log="$git_town_root_dir/errors"

# File path for the log file for the summary
summary_log="$git_town_root_dir/test_log"


# Prints the given text in red.
function echo_failure {

  # Output to terminal
  echo_red $*

  # Output to test log
  echo "$red      $*$normal"  >> $summary_log

  # Output to error log
  echo "$current_SUT $current_context $current_spec_description $1"  >> $failures_log
}


# Prints the given text in green.
function echo_success {

  # Output to terminal
  echo_green $*

  # Output to test log
  echo "$green      $*$normal"  >> $summary_log
}


function reset_logs {
  rm -f $failures_log
  rm -f $summary_log
}

