# frozen_string_literal: true

When(/^(I|my coworker) (?:run|runs|have run) `([^`]+)`$/) do |who, commands|
  user = (who == 'I') ? :developer : :coworker
  in_repository user do
    commands.split(';').each { |command| run command.strip }
  end
end


When(/^I run `([^`]+)` in the "(.+?)" folder$/) do |commands, folder_name|
  in_repository :developer do
    Dir.chdir folder_name
    commands.split(';').each { |command| run command.strip }
  end
end


When(/^I run `(.+?)` and enter an empty commit message$/) do |command|
  # In vim "dG" removes all lines and "ZZ" saves and exits
  @result = run command, inputs: ['dGZZ']
end


When(/^I run `(.+?)` and don't change the default commit message$/) do |command|
  # In vim "ZZ" saves and exits
  @result = run command, inputs: ['ZZ']
end


When(/^I run `(.+?)` and answer the prompts:$/) do |command, table|
  table.map_headers!(&:downcase)
  table.map_column!('answer') do |text|
    text
      .gsub('[ENTER]', "\n")
      .gsub('[DOWN]', "\e[B")
      .gsub('[UP]', "\e[A")
      .gsub('[SPACE]', ' ')
  end
  @result = run command, responses: table.hashes
end




Then(/^it prints the error:$/) do |error_message|
  @error_expected = true
  expect(@last_run_result).to_not be_nil, 'Error message expected, but no commands were run'
  expect(@last_run_result.error).to be_truthy
  actual = unformatted_last_run_output
  expect(actual).to include(error_message), %(
    ACTUAL
    ***************************************************
    #{actual.dump.gsub '\n', "\n"}
    ***************************************************
    EXPECTED TO INCLUDE
    ***************************************************
    #{error_message.dump.gsub '\n', "\n"}
    ***************************************************
  ).gsub(/^ {4}/, '')
end


Then(/^it prints the error "(.+?)"$/) do |error_message|
  step 'it prints the error:', error_message
end


Then(/^it runs no commands$/) do
  expect(commands_of_last_run).to be_empty
end


Then(/^it runs the commands$/) do |expected_commands|
  # We need ERB here to fill in commit SHAs in Git commands
  expected_commands.map_column! 'COMMAND' do |command|
    ERB.new(command).result
  end

  actual_commands = commands_of_last_run(with_branch: expected_commands.column_names.count == 2).table
  expected_commands.diff! actual_commands
end


Then(/^it prints no output$/) do
  expect(@last_run_result.out).to eql ''
end


Then(/^it does not print "(.*)"$/) do |string|
  expect(unformatted_last_run_output).not_to include(string)
end


Then(/^it prints$/) do |string|
  expect(unformatted_last_run_output).to include(string)
end


Then(/^it prints "(.*)"$/) do |string|
  step 'it prints', string
end
