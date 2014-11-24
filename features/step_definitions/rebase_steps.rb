Then(/^my repo (?:still )?has a rebase in progress$/) do
  expect(rebase_in_progress).to be_truthy
end


Then(/^there is no rebase in progress$/) do
  expect(rebase_in_progress).to be_falsy
end
