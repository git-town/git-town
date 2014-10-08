Then /^there is an? (abort|continue) script for "([^"]+)"$/ do |command, operation|
  expect(File.exists? script_path(operation: operation, command: command)).to be_truthy
end


Then /^there is an? (abort|continue) script for "([^"]+)" containing$/ do |command, operation, expected_content|
  file_path = script_path operation: operation, command: command
  expect(File.exists? file_path).to be_truthy
  expect(IO.read file_path).to include expected_content
end


Then /^there is no (abort|continue) script for "(.+?)" anymore$/ do |command, operation|
  expect(File.exists? script_path(operation: operation, command: command)).to be_falsy
end
