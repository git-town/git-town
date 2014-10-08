When /^I successfully finish the rebase by resolving the merge conflict of file "(.*?)"$/ do |file_name|
  IO.write file_name, "resolved content"
  run "git add #{file_name} ; git rebase --continue"
end




Then /^my repo (?:still )?has a rebase in progress$/ do
  expect(rebase_in_progress).to be_truthy
end


Then /^there is no rebase in progress$/ do
  expect(rebase_in_progress).to be_falsy
end
