Given(/^I have a command pattern script$/) do
  mock_script_create
end

Given(/^it has the following preconditions:$/) do |table|
  mock_script_update 'PRECONDITIONS', table.transpose.raw[0].drop(1)
end

Given(/^it has the following steps:$/) do |table|
  mock_script_update 'STEPS', table.transpose.raw[0].drop(1)
end

When(/^I run the command pattern script$/) do
  @result = run 'git mock'
end