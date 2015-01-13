def commits_diff actual, expected
  expected = extract_commits expected
  actual = extract_commits actual

  section_options = [
    ['Expected but not actual commits', expected - actual, skip_if_empty: true],
    ['Actual but not expected commits', actual - expected, skip_if_empty: true]
  ]

  section_options.map { |options| commit_diff_section(*options) }.join('') + "\n"
end


def commit_diff_section title, commits, skip_if_empty: false
  return '' if skip_if_empty && commits.empty?
  "\n#{title}:\n" + commits.map { |c| commit_to_s(c) }.join('')
end


def commit_to_s commit
  "#{commit[:branch]} branch: '#{commit[:message]}' with #{commit[:files]}\n"
end

def extract_commits commit_mapping
  commits = commit_mapping.each_pair.map do |branch, branch_commits|
    branch_commits.map { |commit| commit.merge(branch: branch) }
  end

  commits.flatten.sort_by { |commit| commit_to_s commit }
end
