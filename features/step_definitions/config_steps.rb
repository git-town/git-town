# frozen_string_literal: true

Given(/^I have a the git configuration for "([^"]*)" set to "([^"]*)"$/) do |key, value|
  run_shell_command "git config #{key} #{value}"
end
