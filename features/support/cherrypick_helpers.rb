def cherrypick_in_progress
  get_output('git status').include?("You are currently cherry-picking")
end
