# Returns all files that exist in the current namespace
def existing_files
  Dir['*'] - IGNORED_FILES
end


def uncommitted_files
  get_output_as_array("git status --porcelain | awk '{print $2}'")
end
