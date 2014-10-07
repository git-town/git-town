Given /^I am on a feature branch$/ do
  run "git checkout -b feature main"
  run "git push -u origin feature"
end


Given /^I am on a local feature branch$/ do
  run "git checkout -b feature main"
end


Given /^I am on the main branch$/ do
  run "git checkout main"
end


Given /^I have a feature branch named "(.*)"$/ do |branch_name|
  run "git branch #{branch_name} main"
end





Then /^I (?:end up|am still) on the feature branch$/ do
  expect(current_branch_name).to eql 'feature'
end


Then /^I (?:end up|am still) on the main branch$/ do
  expect(current_branch_name).to eql 'main'
end


Then /^I end up on my feature branch$/  do
  expect(current_branch_name).to eql 'feature'
end


Then /^I (?:end up|am still) on the "(.+?)" branch$/ do |branch_name|
  expect(current_branch_name).to eql branch_name
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
