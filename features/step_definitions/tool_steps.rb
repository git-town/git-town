Given(/^I have (.+?) installed$/) do |tool|
  tools = (tool == 'no tools') ? [] : [tool]
  update_installed_tools tools
end
