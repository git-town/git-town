Given /^I am on a feature branch$/ do
  create_branch "feature"
end

Given /^I am on a local feature branch$/ do
  run "git checkout -b feature main"
end


Given /^I am on the main branch$/ do
  run "git checkout main"
end


Given /^I am on the "(.+?)" branch$/ do |branch_name|
  if existing_local_branches.include?(branch_name)
    run "git checkout #{branch_name}"
  else
    create_branch branch_name
  end
end


Given /^I have a feature branch named "(.+?)"(?: (behind|ahead of) main)?$/ do |branch_name, relation|
  create_branch branch_name
  if relation
    commit_to_branch = relation == 'behind' ? 'main' : branch_name
    create_commits branch: commit_to_branch
  end
end


Given /^I have a non\-feature branch "(.+?)" behind main$/ do |branch_name|
  create_branch branch_name
  set_non_feature_branches_configuration branch_name
  create_commits branch: 'main'
end


Given /^my coworker has a feature branch named "(.*)"(?: (behind|ahead of) main)?$/ do |branch_name, relation|
  at_path coworker_repository_path do
    create_branch branch_name
    if relation
      commit_to_branch = relation == 'behind' ? 'main' : branch_name
      create_commits branch: commit_to_branch
    end
  end
end




Then /^I (?:end up|am still) on the feature branch$/ do
  expect(current_branch_name).to eql 'feature'
end


Then /^I end up on my feature branch$/ do
  expect(current_branch_name).to eql 'feature'
end


Then /^I (?:end up|am still) on the "(.+?)" branch$/ do |branch_name|
  expect(current_branch_name).to eql branch_name
end


Then /^there is no "(.+?)" branch$/ do |branch_name|
  expect(existing_local_branches).to_not include(branch_name)
end


Then /^my coworker ends up on the "(.+?)" branch$/ do |branch_name|
  at_path coworker_repository_path do
    expect(current_branch_name).to eql branch_name
  end
end


Then /^the branch "(.*?)" has not been pushed to the repository$/ do |branch_name|
  expect(existing_remote_branches).to_not include(branch_name)
end


Then /^all branches are now synchronized$/ do
  expect(number_of_branches_out_of_sync).to eql 0
end


Then /^there are no more feature branches$/ do
  expected_branches = [ 'master',
                        "* main",
                        "remotes/origin/main",
                        'remotes/origin/master' ].sort
  actual_branches = run("git branch -a").out.split("\n").map(&:strip).sort
  expect(actual_branches).to eql expected_branches
end


Then /^the branch "(.*?)" is deleted everywhere$/ do |branch_name|
  expect(existing_local_branches).to_not include(branch_name)

  at_path coworker_repository_path do
    expect(existing_local_branches).to_not include(branch_name)
  end

  expect(existing_remote_branches).to_not include("remotes/origin/#{branch_name}")
end


Then /^the branch "(.+?)" still exists$/ do |branch_name|
  expect(existing_remote_branches).to include("remotes/origin/#{branch_name}")
end


Then /^the existing branches are$/ do |table|
  verify_branches table.hashes
end

