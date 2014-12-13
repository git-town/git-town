Given(/^I don't have a main branch name configured$/) do
  delete_main_branch_configuration
end


Given(/^my non\-feature branches are not configured$/) do
  delete_non_feature_branches_configuration
end


Given(/^I have an old configuration file with (.+?)$/) do |data|
  delete_main_branch_configuration
  data = Kappamaki.attributes_from_sentence data
  IO.write '.main_branch_name', data[:'main branch']
end


Given(/^non-feature branch configuration "(.+?)"$/) do |configuration|
  configure_non_feature_branches configuration
  configuration.split(',').map(&:strip).each { |b| create_branch b }
end


Given(/^I have configured the main branch name as "(.*)"$/) do |main_branch_name|
  set_configuration 'main-branch-name', main_branch_name
end


Given(/^my non-feature branch(?:es are| is) "(.*)"$/) do |data|
  non_feature_branches = Kappamaki.from_sentence(data).join(", ")
  set_configuration 'non-feature-branch-names', non_feature_branches
end


Given(/I haven't configured Git Town yet/) do
  delete_main_branch_configuration
  delete_non_feature_branches_configuration
end


Then(/^I don't have an old configuration file anymore$/) do
  expect(File.exist? '.main_branch_name').to be_falsy
end


Then(/^the main branch name is now configured as "(.+?)"$/) do |main_branch_name|
  expect(main_branch_configuration).to eql main_branch_name
end


Then(/^the non\-feature branches( don't)? include "(.*?)"$/) do |negate, non_feature_branch|
  if negate
    expect(non_feature_branch_configuration).not_to include non_feature_branch
  else
    expect(non_feature_branch_configuration).to include non_feature_branch
  end
end
