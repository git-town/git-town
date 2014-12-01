def commits_diff actual, expected
  expected.sort_by! { |c| commit_to_s(c) }
  actual.sort_by! { |c| commit_to_s(c) }

  section_options = [
    ['Expected commits', expected],
    ['Actual commits', actual],
    ['Common commits', expected & actual, skip_if_empty: true],
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
