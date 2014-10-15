# Helper methods for dealing with files and temp files.


# Unique string that identifies the current directory and git branch
temp_filename_suffix="_`pwd | tr '/' '_' `"

# Path to the temp file used by these scripts.
user_input_filename="/tmp/git-town-user-input$temp_filename_suffix"


# Removes the temp file.
function delete_temp_file {
  rm $user_input_filename
}


# Ensures that the given tool is installed.
function ensure_tool_installed {
  local toolname=$1
  if [ `which $toolname | wc -l` == 0 ]; then
    echo_error_header
    echo "  You need the '$toolname' tool in order to run tests."
    echo "  Please install it using your package manager,"
    echo "  or on OS X with 'brew install $toolname'."
    exit_with_error
  fi
}
