# Pops the uncommitted changes back from the stash
def restore_open_changes
  run "git stash pop"
end


# Stashes away uncommitted changes
def stash_open_changes
  run "git stash -u"
end


# Returns the number of unpopped stashes
def stash_size
  integer_output_of "git stash list | wc -l"
end


# Executes the given block with the uncommitted changes stashed away
def without_open_changes
  do_stash = has_uncommitted_files
  stash_open_changes if do_stash
  yield
  restore_open_changes if do_stash
end
