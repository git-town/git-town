Given(/^I have "(.+?)" installed$/) do |tool|
  update_installed_tools [tool]
end


Given(/^I have no command that opens browsers installed$/) do
  update_installed_tools []
end
