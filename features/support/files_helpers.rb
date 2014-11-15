# Returns all files that exist in the current namespace
def existing_files
  Dir['*'] - IGNORED_FILES
end


# Returns whether the current workspace has uncommitted files
def has_uncommitted_files
  uncommitted_files.size > 0
end


def uncommitted_files
  array_output_of "git status --porcelain | awk '{print $2}'"
end
