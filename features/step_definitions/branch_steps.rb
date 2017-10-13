# frozen_string_literal: true

Given(/^(I|my coworker) (?:am|is) on the "(.+?)" branch$/) do |who, branch_name|
  user = (who == 'I') ? :developer : :coworker
  in_repository user do
    run "git checkout #{branch_name}"
  end
end


Given(/^my repository has a( local)?( feature)?( perennial)? branch named "([^"]+)"( on another machine)?$/) do |local, feature, perennial, branch_name, remote|
  user = 'developer'
  user += '_secondary' if remote
  in_repository user do
    create_branch branch_name, remote: !local
    set_parent_branch branch: branch_name, parent: 'main', ancestors: 'main' if feature
    add_perennial_branch branch_name if perennial
  end
end


Given(/^my repository has a feature branch named "([^"]+)" with no parent$/) do |branch_name|
  create_branch branch_name
end


Given(/^my repository has the( local)?( feature)?( perennial)? branches "(.+?)"$/) do |local, feature, perennial, branch_names|
  Kappamaki.from_sentence(branch_names).each do |branch_name|
    create_branch branch_name, remote: !local
    set_parent_branch branch: branch_name, parent: 'main', ancestors: 'main' if feature
    add_perennial_branch branch_name if perennial
  end
end


Given(/^(?:my repository|it) has a(?: feature| hotfix)? branch named "([^"]+)" as a child of "([^"]+)"$/) do |branch_name, parent_name|
  create_branch branch_name, remote: true, start_point: parent_name
  set_parent_branch branch: branch_name, parent: parent_name
  store_branch_hierarchy_metadata
end


Given(/^my coworker has a feature branch named "(.+?)"(?: (behind|ahead of) main)?$/) do |branch_name, relation|
  in_repository :coworker do
    create_branch branch_name
    if relation
      commit_to_branch = relation == 'behind' ? 'main' : branch_name
      create_commits branch: commit_to_branch
    end
  end
end


Given(/^my repository knows about the remote branch$/) do
  run 'git fetch'
end


Given(/the "(.+?)" branch gets deleted on the remote/) do |branch_name|
  in_repository :coworker do
    run "git push origin :#{branch_name}"
  end
end


Given(/^I am on the "(.+?)" branch with "(.+?)" as the previous Git branch/) do |current_branch, previous_branch|
  run "git checkout #{previous_branch}"
  run "git checkout #{current_branch}"
end


Given(/^(I|my coworker) sets? the parent branch of "([^"]*)" as "([^"]*)"$/) do |who, child_branch, parent_branch|
  user = (who == 'I') ? :developer : :coworker
  in_repository user do
    set_parent_branch branch: child_branch, parent: parent_branch
  end
end



Then(/^my previous Git branch is (?:now|still) "(.+?)"/) do |previous_branch|
  run 'git checkout -'
  expect(current_branch_name).to eql previous_branch
  run 'git checkout -'
end


Then(/^I (?:end up|am still) on the "(.+?)" branch$/) do |branch_name|
  expect(current_branch_name).to eql branch_name
end


Then(/^there is no "(.+?)" branch$/) do |branch_name|
  expect(existing_local_branches).to_not include(branch_name)
  expect(existing_remote_branches).to_not include("origin/#{branch_name}")
end


Then(/^all branches are now synchronized$/) do
  expect(number_of_branches_out_of_sync).to eql 0
end


Then(/^there are no more feature branches$/) do
  expected_branches = ['main', 'origin/main']
  perennial_branches.each do |perennial_branch|
    expected_branches << perennial_branch
    expected_branches << "origin/#{perennial_branch}"
  end
  expect(existing_branches).to match_array expected_branches
end


Then(/^the existing branches are$/) do |table|
  verify_branches table
end
