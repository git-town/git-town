def cherrypick_in_progress
  run('git status | grep "You are currently cherry-picking" | wc -l')[:out] == '1'
end
