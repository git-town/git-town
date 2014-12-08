Given(/^I add a local tag "(.+?)"$/) do |tag_name|
  run "git tag -a #{tag_name} -m '#{tag_name}'"
end




Then(/^tag "(.+?)" has been pushed to the remote$/) do |tag_name|
  expect(unpushed_tags).to_not include tag_name
end
