# Helper methods for writing to the terminal.


function echo_all_done {
  echo
  echo_header "*** ALL DONE ***"
  echo
}


# Prints a header line into the terminal.
function echo_header {
  echo
  tput bold
  echo "$*"
  tput sgr0
}


# Prints an error header into the terminal.
function echo_error_header {
  echo
  tput bold
  echo "  Error"
  tput sgr0
}


# Prints the intro line of a script into the terminal.
function echo_intro {
  tput bold
  echo "$*"
  tput sgr0
}


# Outputs the given text in red
function echo_red {
  tput setaf 1
  echo "$*"
  tput sgr0
}


# Prints the header for the usage help screen into the terminal.
function echo_usage_header {
  echo
  tput bold
  echo "  Usage"
  tput sgr0
}


# Exits the currently running script with an error response code.
function exit_with_error {
  echo
  echo
  exit 1
}
