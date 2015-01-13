Given(/^the following commits? exists? in (my|my coworker's) repository$/) do |who, commits_table|
  user = (who == 'my') ? :developer : :coworker
  commits_table.map_headers!(&:downcase)
  @initial_commits_table = commits_table.clone
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


Then(/^I am left with my original commits$/) do
  verify_commits @initial_commits_table.hashes
end


Then(/^there are no commits$/) do
  expect(commits_in_repo).to be_empty
end
