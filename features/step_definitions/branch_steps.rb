Given(/^(I|my coworker) (?:am|is) on the "(.+?)" branch$/) do |who, branch_name|
  user = (who == 'I') ? :developer : :coworker
  in_repository user do
    run "git checkout #{branch_name}"
  end
end


Given(/^I have(?: a)?( local)?(?: feature)? branch(?:es)? named "(.+?)"$/) do |local, branch_names|
  Kappamaki.from_sentence(branch_names).each do |branch_name|
    create_branch branch_name, remote: !local
    set_parent_branch branch: branch_name, parent: 'main', parents: 'main'
  end
end


Given(/^I have a( local)? feature branch named "(.+?)" (behind|ahead of) main$/) do |local, branch_name, relation|
  create_branch branch_name, remote: !local
  if relation
    commit_to_branch = relation == 'behind' ? 'main' : branch_name
    create_commits branch: commit_to_branch
  end
end


Given(/^I have a feature branch named "(.+?)" that is cut from (.+?)$/) do |branch_name, parent_branch_name|
  create_branch branch_name, start_point: parent_branch_name
end


Given(/^I have a stale feature branch named "([^"]+)" with its tip at "([^"]+)"$/) do |branch_name, commit_message|
  create_branch branch_name, start_point: commit_sha(commit_message)
end


Given(/^I have a non\-feature branch "(.+?)" behind main$/) do |branch_name|
  create_branch branch_name
  configure_non_feature_branches branch_name
  create_commits branch: 'main'
end


Given(/^Git Town knows that "(.+?)" has the parent "(.+?)" and the parents "(.+?)"$/) do |branch, parent, parents|
  set_parent_branch branch: branch, parent: parent, parents: parents
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


Given(/the "(.+?)" branch gets deleted on the remote/) do |branch_name|
  in_repository :coworker do
    run "git push origin :#{branch_name}"
  end
end





Then(/^I (?:end up|am still) on the "(.+?)" branch$/) do |branch_name|
  expect(current_branch_name).to eql branch_name
end


Then(/^there is no "(.+?)" branch$/) do |branch_name|
  expect(existing_local_branches).to_not include(branch_name)
  expect(existing_remote_branches).to_not include("origin/#{branch_name}")
end


Then(/^the branch "(.+?)" has not been pushed to the repository$/) do |branch_name|
  expect(existing_remote_branches).to_not include(branch_name)
end


Then(/^all branches are now synchronized$/) do
  expect(number_of_branches_out_of_sync).to eql 0
end


Then(/^there are no more feature branches$/) do
  expect(existing_branches).to match_array ['main', 'origin/main']
end


Then(/^the existing branches are$/) do |table|
  verify_branches table
end
