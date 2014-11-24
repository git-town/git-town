Then(/^my repo(?: still)? has a merge in progress$/) do
  expect(merge_in_progress?).to be_truthy
end


Then(/^there is no merge in progress$/) do
  expect(merge_in_progress?).to be_falsy
end
