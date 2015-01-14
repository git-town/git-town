def commits_diff actual, expected
  out = ''

  (expected.keys + actual.keys).uniq.sort.each do |branch|
    next if expected[branch] == actual[branch]
    out += "\n#{branch} branch\n"
    out += commit_list_with_title 'Expected commits', expected[branch]
    out += commit_list_with_title 'Actual commits', actual[branch]
  end

  out
end


def commit_to_s commit
  "    '#{commit[:message]}' with #{commit[:files]}\n"
end


def commit_list commits
  commits.map { |commit| commit_to_s commit }.join('')
end


def commit_list_with_title title, commits
  "  #{title}\n#{commit_list commits}"
end
