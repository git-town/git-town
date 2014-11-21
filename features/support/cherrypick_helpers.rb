def cherrypick_in_progress
  output_of('git status').include? 'You are currently cherry-picking'
end
