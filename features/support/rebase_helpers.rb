def rebase_in_progress
  output_of('git status').include? 'You are currently rebasing'
end
