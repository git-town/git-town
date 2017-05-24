# frozen_string_literal: true
# Returns the number of unpopped stashes
def stash_size
  integer_output_of 'git stash list | wc -l'
end
