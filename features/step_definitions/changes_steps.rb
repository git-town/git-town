Then /^there are no open changes$/ do
  expect(run_this("git status --short")[:out]).to eql ''
end

