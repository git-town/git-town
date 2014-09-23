# Returns all files that exist in the current namespace
def existing_files
  Dir['*'] - IGNORED_FILES
end
