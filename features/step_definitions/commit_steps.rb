Given(/^the following commits? exists? in (my|my coworker's) repository$/) do |who, commits_table|
  user = (who == 'my') ? :developer : :coworker
  commits_table.map_headers!(&:downcase)
  in_repository user do
    create_commits commits_table.hashes
  end
end




Then(/^(?:now )?(I|my coworker) (?:still )?(?:have|has) the following commits$/) do |who, commits_table|
  user = (who == 'I') ? :developer : :coworker
  commits_table.map_headers!(&:downcase)
  in_repository user do
    verify_commits commits_table.hashes
  end
end


Then(/^there are no commits$/) do
  expect(commits_in_repo).to eql []
end
