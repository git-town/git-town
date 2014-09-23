Then /^my repo has a cherry\-pick in progress$/ do
  expect(cherrypick_in_progress).to be_truthy
end


Then /^my repo has no cherry\-pick in progress$/ do
  expect(cherrypick_in_progress).to be_falsy
end

