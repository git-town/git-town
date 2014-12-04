Given(/^the following commits? exists? in my repository$/) do |commits_table|
  at_path local_repository_path do
    create_commits commits_table.hashes
  end
end


Given(/^the following commits? exists? in Charlie's repository$/) do |commits_table|
  at_path coworker_repository_path do
    create_commits commits_table.hashes
  end
end




Then(/^(?:now )?(?:(?:I (?:still )?(?:have|see))) the following commits$/) do |commits_table|
  verify_commits commits_table: commits_table, repository_path: local_repository_path
end


Then(/^(?:now )?Charlie(?: still)? sees the following commits$/) do |commits_table|
  verify_commits commits_table: commits_table, repository_path: coworker_repository_path
end


Then(/^there are no commits$/) do
  expect(commits_in_repo).to eql []
end
