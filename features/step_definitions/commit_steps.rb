Given(/^the following commits? exists? in (my|my coworker's) repository$/) do |who, commits_table|
  path = (who == 'my') ? local_repository_path : coworker_repository_path
  commits_table.map_headers!(&:downcase)
  @initial_commits_table = commits_table.clone
  at_path(path) do
    create_commits commits_table.hashes
  end
end




Then(/^(?:now )?(I|my coworker) (?:still )?(?:have|has) the following commits$/) do |who, commits_table|
  path = (who == 'I') ? local_repository_path : coworker_repository_path
  commits_table.map_headers!(&:downcase)
  at_path(path) do
    verify_commits commits_table.hashes
  end
end


Then(/^I am left with my original commits$/) do
  at_path(local_repository_path) do
    verify_commits @initial_commits_table.hashes
  end
end


Then(/^there are no commits$/) do
  expect(commits_in_repo).to eql []
end
