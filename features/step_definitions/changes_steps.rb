Then /^there are no open changes$/ do
  expect(run("git status --short").out).to eql ''
end

