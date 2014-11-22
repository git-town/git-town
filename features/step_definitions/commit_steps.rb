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




Then(/^my branch and its remote still have (\d+) and (\d+) different commits$/) do |local_count, remote_count|
  matches = /have (\d+) and (\d+) different commit each/.match(run("git status").out)
  expect(matches[1]).to eql local_count
  expect(matches[2]).to eql remote_count
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
