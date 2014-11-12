# Returns the number of unpopped stashes
def stash_size
  get_output_as_integer("git stash list | wc -l")
end

