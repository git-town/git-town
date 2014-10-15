# Helper methods for writing to the terminal.


# Prints a line in bold
function echo_bold {
  tput bold
  echo "$*"
  tput sgr0
}


# Prints a header line into the terminal.
function echo_header {
  echo
  echo_bold "$*"
}


# Prints an error header into the terminal.
function echo_error_header {
  echo_header "  Error"
}


# Outputs the given text in red
function echo_red {
  tput setaf 1
  echo "$*"
  tput sgr0
}


# Prints the header for the usage help screen into the terminal.
function echo_usage_header {
  echo_header "  Usage"
}


# Exits the currently running script with an error response code.
function exit_with_error {
  echo
  exit 1
}


# Exits the currently running script with a success response code
function exit_with_success {
  echo
  exit 0
}


# Prints a command
function print_command {
  echo
  echo_bold $*
  commands_printed=$((commands_printed+1))
}


# Run a command, also prints command and output
command_exit_status=0
function run_command {
  local cmd=$*
  local output; output=`$cmd 2>&1`
  command_exit_status=$?

  print_command $cmd
  if [ -n "$output" ]; then
    echo "$output"
  fi
}

