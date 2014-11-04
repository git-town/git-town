# Returns the number of unpopped stashes
def stash_size
  run("git stash list | wc -l").out.to_i
end

