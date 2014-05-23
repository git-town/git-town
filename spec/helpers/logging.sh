# Collects error messages and information about which tests pass and which fail.

# File path of the log file for errors
failures_log="$git_town_root_dir/errors"

# File path for the log file for the summary
summary_log="$git_town_root_dir/test_log"


# Prints the given text in red.
function echo_failure {

  # Output to terminal
  tput setaf 1
  echo $*
  tput sgr0

  # Output to test log
  tput setaf 1      >> $summary_log
  echo "      $*"   >> $summary_log
  tput sgr0         >> $summary_log

  # Output to error log
  echo "$current_SUT $current_context $current_spec_description $1"  >> $failures_log
}


# Prints the given text in green.
function echo_success {
  tput setaf 2
  echo $*
  tput sgr0

  # Output to test log
  tput setaf 2      >> $summary_log
  echo "      $*"   >> $summary_log
  tput sgr0         >> $summary_log
}


function reset_logs {
  rm -f $failures_log
  rm -f $summary_log
}

