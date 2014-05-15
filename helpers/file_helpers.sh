# Helper methods for dealing with files and temp files.


# Path to the temp file used by these scripts.
temp_filename="/tmp/git_input_$$"


# Removes the temp file.
function delete_temp_file {
  rm $temp_filename
}

