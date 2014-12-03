Given(/^I don't have a main branch name configured$/) do
  delete_main_branch_configuration
end


Given(/^I have an old configuration file with (.*?)$/) do |data|
  delete_main_branch_configuration
  data = Kappamaki.attributes_from_sentence data
  IO.write '.main_branch_name', data[:'main branch']
end


Given(/^non-feature branch configuration "(.*)"$/) do |configuration|
  configure_non_feature_branches configuration
  configuration.split(',').map(&:strip).each { |b| create_branch b }
end


Given(/^I have set "(.*)" to "(.*)"$/) do |config, value|
  set_configuration config, value
end


Given(/^I am using Git Town version "(.*)"/) do |version|
  set_environment_variable 'GIT_TOWN_VERSION', version
end


Then(/^I don't have an old configuration file anymore$/) do
  expect(File.exist? '.main_branch_name').to be_falsy
end


Then(/^the main branch name is now configured as "(.*?)"$/) do |main_branch_name|
  expect(main_branch_configuration).to eql main_branch_name
end
