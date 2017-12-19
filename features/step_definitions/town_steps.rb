# frozen_string_literal: true

Given(/^Git Town is in offline mode$/) do
  set_global_configuration 'offline', true
end


Given(/^Git Town is not in offline mode$/) do
  set_global_configuration 'offline', false
end


Given(/^I don't have a main branch name configured$/) do
  delete_main_branch_configuration
end


Given(/^my perennial branches are not configured$/) do
  delete_perennial_branches_configuration
end


Given(/^the "([^"]*)" configuration is set to "([^"]*)"$/) do |configuration, value|
  set_configuration configuration, value
end


Given(/^the global "([^"]*)" configuration is set to "([^"]*)"$/) do |configuration, value|
  set_global_configuration configuration, value
end


Given(/^the main branch is configured as "(.*)"$/) do |main_branch_name|
  set_configuration 'main-branch-name', main_branch_name
end


Given(/^the perennial branches are configured as (.*)$/) do |data|
  branch_names = Kappamaki.from_sentence data
  set_configuration 'perennial-branch-names', branch_names.join(' ')
end


Given(/I haven't configured Git Town yet/) do
  delete_main_branch_configuration
  delete_perennial_branches_configuration
end


Given(/^I configure "([^"]*)" as "([^"]*)"$/) do |key, value|
  run "git config '#{key}' '#{value}'"
end



Then(/^the perennial branches are now configured as "(.*)"$/) do |data|
  branch_names = Kappamaki.from_sentence(data)
  expect(perennial_branch_configuration.split(' ').map(&:strip)).to match_array branch_names
end


Then(/^my repo is configured with no perennial branches$/) do
  expect(perennial_branch_configuration).to be_empty
end


Then(/^the main branch is now configured as "([^"]*)"$/) do |branch_name|
  expect(main_branch_configuration).to eql branch_name
end


Then(/^my repo is now configured with "([^"]*)" set to "(.+?)"$/) do |configuration, value|
  expect(get_configuration(configuration)).to eql value
end


Then(/^git is now configured with "([^"]*)" set to "(.+?)"$/) do |configuration, value|
  expect(get_global_configuration(configuration)).to eql value
end


Then(/^Git Town is (?:no longer|still not) configured for this repository$/) do
  expect(git_town_configuration).to be_empty
end


Then(/^it prints the initial configuration prompt$/) do
  step %(it prints "Git Town needs to be configured")
end


Then(/^offline mode is enabled$/) do
  expect(get_configuration('offline')).to eql 'true'
end


Then(/^offline mode is disabled$/) do
  expect(get_configuration('offline')).to eql 'false'
end
