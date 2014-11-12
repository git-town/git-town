def merge_in_progress?
  get_output('git status').include?('You have unmerged paths')
end
