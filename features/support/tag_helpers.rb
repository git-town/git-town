# Returns the names of all unpushed tags.
def unpushed_tags
  output_of("git push --tags --dry-run 2>&1 | grep 'new tag' | awk '{print $4}'")
end
