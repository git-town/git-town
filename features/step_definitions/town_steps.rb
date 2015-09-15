Given(/^I don't have a main branch name configured$/) do
  delete_main_branch_configuration
end


Given(/^my perennial branches are not configured$/) do
  delete_perennial_branches_configuration
end


Given(/^I have an old configuration file with (.+?)$/) do |data|
  delete_main_branch_configuration
  data = Kappamaki.attributes_from_sentence data
  IO.write '.main_branch_name', data[:'main branch']
end


Given(/^I have configured the main branch name as "(.*)"$/) do |main_branch_name|
  set_configuration 'main-branch-name', main_branch_name
end


Given(/^my non-feature branches are configured as (.*)$/) do |data|
  branch_names = Kappamaki.from_sentence data
  set_configuration 'non-feature-branch-names', branch_names.join(' ')
end


Given(/^my perennial branches are configured as (.*)$/) do |data|
  branch_names = Kappamaki.from_sentence data
  set_configuration 'perennial-branch-names', branch_names.join(' ')
end


Given(/I haven't configured Git Town yet/) do
  delete_main_branch_configuration
  delete_perennial_branches_configuration
end




Then(/^I don't have an old configuration file anymore$/) do
  expect(File.exist? '.main_branch_name').to be_falsy
end


Then(/^I have no non\-feature branch configuration$/) do
  expect(non_feature_branch_configuration).to eql ''
end


Then(/^my repo is configured with perennial branches as "(.*)"$/) do |data|
  branch_names = Kappamaki.from_sentence(data)
  expect(perennial_branch_configuration.split(' ').map(&:strip)).to match_array branch_names
end


Then(/^my repo is configured with no perennial branches$/) do
  expect(perennial_branch_configuration).to be_empty
end


Then(/^my repo is configured with the main branch as "([^"]*)"$/) do |branch_name|
  expect(main_branch_configuration).to eql branch_name
end


Then(/^Git Town is (?:no longer|still not) configured for this repository$/) do
  expect(git_town_configuration).to be_empty
end


Then(/^I see the initial configuration prompt$/) do
  step %(I see "Git Town needs to be configured")
end


Then(/^I (don't )?see the first line of the configuration wizard$/) do |negate|
  configuration_wizard_first_line = 'Please specify the main dev branch'
  if negate
    step %(I don't see "#{configuration_wizard_first_line}")
  else
    step %(I see "#{configuration_wizard_first_line}")
  end
end
