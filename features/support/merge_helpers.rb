def merge_in_progress?
  run("git status | grep 'You have unmerged paths' | wc -l")[:out] == '1'
end
