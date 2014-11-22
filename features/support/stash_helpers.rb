# Pops the uncommitted changes back from the stash
def restore_open_changes
  run "git stash pop"
end


# Returns the number of unpopped stashes
def stash_size
  integer_output_of 'git stash list | wc -l'
end
