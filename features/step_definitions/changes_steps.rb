# frozen_string_literal: true

Then(/^there are no open changes$/) do
  changes = run('git status --short').out
  changes.sub!(/.* coverage.cov\n/, '')
  expect(changes).to eql ''
end
