Then /^there are abort and continue scripts for "(.+?)"$/  do |operation|
  expect(File.exists? script_path(operation: operation, command: 'abort')).to be_truthy
  expect(File.exists? script_path(operation: operation, command: 'continue')).to be_truthy
end


Then /^there are no abort and continue scripts for "(.*?)" anymore$/ do |operation|
  expect(File.exists? script_path(operation: operation, command: 'abort')).to be_falsy
  expect(File.exists? script_path(operation: operation, command: 'continue')).to be_falsy
end


Then /^there is an? (abort|continue) script for "(.+?)"$/ do |command, operation|
  expect(File.exists? script_path(operation: operation, command: command)).to be_truthy
end


Then /^there is no (abort|continue) script for "(.+?)" anymore$/ do |command, operation|
  expect(File.exists? script_path(operation: operation, command: command)).to be_falsy
end
