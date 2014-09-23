def rebase_in_progress
  run_this("git status | grep 'You are currently rebasing' | wc -l")[:out] == '1'
end
