When(/^(I|my coworker) runs? `([^`]+)`$/) do |who, commands|
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


When(/^I run `([^`]+)` with the last( two)? commit shas?$/) do |command, two|
  count = two ? 2 : 1
  shas = recent_commit_shas(count).join(' ')
  step "I run `#{command} #{shas}`"
end


When(/^I run `(.+?)` and enter "(.*?)"$/) do |command, input|
  inputs = Kappamaki.from_sentence(input)
  @result = run command, inputs: inputs
end


When(/^I run `(.+?)` and enter an empty commit message$/) do |command|
  # In vim "dG" removes all lines and "ZZ" saves and exits
  step "I run `#{command}` and enter \"dGZZ\""
end


When(/^I run `(.+?)` and enter main branch name "(.+?)"(?: and non\-feature branch names "(.+)")?/) do |cmd, main, non_feature|
  @result = run cmd, inputs: [main, non_feature].compact
end




Then(/^I get the error$/) do |error_message|
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


Then(/^I get the error "(.+?)"$/) do |error_message|
  step 'I get the error', error_message
end


Then(/^it runs no (?:Git|shell) commands$/) do
  expect(commands_of_last_run).to be_empty
end


Then(/^it runs the Git commands$/) do |expected_commands|
  # We need ERB here to fill in commit SHAs in Git commands
  expected_commands.map_column! 'COMMAND' do |command|
    ERB.new(command).result
  end

  expected_commands.diff! commands_of_last_run.table
end


Then(/^it runs the following shell commands/) do |expected_commands|
  expected_commands.map_column! 'COMMAND' do |command|
    ERB.new(command).result
  end
  expected_commands.diff! commands_of_last_run_outside_git.unshift(expected_commands.headers)
end


Then(/^I see no output$/) do
  expect(@last_run_result.out).to eql ''
end


Then(/^I don't see "(.*)"$/) do |string|
  expect(@last_run_result.out).not_to include(string)
end


Then(/^I see$/) do |string|
  expect(unformatted_last_run_output).to include(string)
end


Then(/^I see "(.*)"$/) do |string|
  step 'I see', string
end


Then(/^I see the "(.+?)" man page$/) do |manpage|
  expect(@last_run_result.out).to eql "man called with: #{manpage}\n"
end
