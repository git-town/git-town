def merge_in_progress?
  output_of('git status').include? 'You have unmerged paths'
end
