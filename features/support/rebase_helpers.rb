def rebase_in_progress
  get_output('git status').include?('You are currently rebasing')
end
