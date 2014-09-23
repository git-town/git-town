Then /^there is an? (abort|continue) script for "(.+?)"$/ do |command, operation|
  expect(File.exists? script_path(operation: operation, command: command)).to be_truthy
end


Then /^there is no (abort|continue) script for "(.+?)" anymore$/ do |command, operation|
  expect(File.exists? script_path(operation: operation, command: command)).to be_falsy
end
