# Returns all files that exist in the current namespace
def existing_files
  Dir['*'] - IGNORED_FILES
end


def uncommitted_files
  array_output_of "git status --porcelain | awk '{print $2}'"
end


# Runs the given block with a file with the given name and content
def with_file name, content
  IO.write name, content
  yield
  File.delete name
end
