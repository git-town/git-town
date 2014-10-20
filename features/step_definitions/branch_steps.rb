Given /^I am on a feature branch$/ do
  create_branch "feature", checkout: true
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
    create_branch branch_name, checkout: true
  end
end


Given /^I have a(?: feature)? branch named "(.*)"$/ do |branch_name|
  create_branch branch_name, checkout: false
end


Given /^my coworker has a feature branch named "(.*)"$/ do |branch_name|
  at_path coworker_repository_path do
    create_branch branch_name, checkout: false
  end
end



Then /^I (?:end up|am still) on the feature branch$/ do
  expect(current_branch_name).to eql 'feature'
end


Then /^I end up on my feature branch$/  do
  expect(current_branch_name).to eql 'feature'
end


Then /^I (?:end up|am still) on the "(.+?)" branch$/ do |branch_name|
  expect(current_branch_name).to eql branch_name
end


Then /^I have the feature branches (.+?)$/ do |branch_names|
  expect(existing_feature_branches).to eq Kappamaki.from_sentence(branch_names)
end


Then /^the branch "(.*?)" has not been pushed to the repository$/ do |branch_name|
  expect(remote_branch_exists branch_name).to be_falsy
end


Then /^all branches are now synchronized$/ do
  run("git branch -vv | grep $1 | grep -o '\[.*\]' | tr -d '[]' | awk '{ print $2 }' | tr -d '\n' | wc -m") == '0'
end


Then /^there are no more feature branches$/ do
  expected_branches = [ 'master',
                        "* main",
                        "remotes/origin/main",
                        'remotes/origin/master' ].sort
  actual_branches = run("git branch -a")[:out].split("\n").map(&:strip).sort
  expect(actual_branches).to eql expected_branches
end

Then /^the branch "(.+?)" will be deleted$/  do |branch_name|
  expect(existing_local_branches).to_not include(branch_name)
  expect(existing_remote_branches).to_not include(branch_name)
end

Then(/^the branch "(.+?)" still exists$/) do |branch_name|
  expect(existing_remote_branches).to include("remotes/origin/#{branch_name}")
end

