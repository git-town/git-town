When(/^I successfully finish the merge by resolving the conflict in "(.*?)"$/) do |file_name|
  IO.write file_name, 'resolved content'
  run "git add #{file_name} ; git commit -m \"Merge branch 'main' into feature\""
end




Then(/^my repo(?: still)? has a merge in progress$/) do
  expect(merge_in_progress?).to be_truthy
end


Then(/^there is no merge in progress$/) do
  expect(merge_in_progress?).to be_falsy
end
