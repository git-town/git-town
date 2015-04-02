Given(/^the following commits? exists? in (my|my coworker's) repository$/) do |who, commits_table|
  user = (who == 'my') ? :developer : :coworker
  commits_table.map_headers!(&:downcase)
  @initial_commits_table = commits_table.clone
  in_repository user do
    create_commits commits_table.hashes
    @original_files = files_in_branches
  end
end




Then(/^(?:now )?(I|my coworker) (?:still )?(?:have|has) the following commits$/) do |who, commits_table|
  user = (who == 'I') ? :developer : :coworker
  in_repository user do
    verify_commits commits_table
  end
end


Then(/^I am left with my original commits$/) do
  @initial_commits_table.map_headers!(&:upcase)
  verify_commits @initial_commits_table
end


Then(/^there are no commits$/) do
  expect(CommitListBuilder.new.add_commits_in_current_repo).to be_empty
end
