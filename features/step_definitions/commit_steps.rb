Given(/^the following commits? exists? in (my|my coworker's) repository$/) do |who, commits_table|
  path = (who == 'my') ? local_repository_path : coworker_repository_path
  commits_table.map_headers!(&:downcase)
  at_path(path) { create_commits commits_table.hashes }
end




Then(/^(?:now )?(I|my coworker) (?:still )?(?:have|has) the following commits$/) do |who, commits_table|
  path = (who == 'I') ? local_repository_path : coworker_repository_path
  commits_table.map_headers!(&:downcase)
  at_path(path) { verify_commits commits_table.hashes }
end


Then(/^there are no commits$/) do
  expect(commits_in_repo).to eql []
end
