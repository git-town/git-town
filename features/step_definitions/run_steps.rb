When(/^(I|my coworker) runs? `([^`]+)`$/) do |who, commands|
  user = (who == 'I') ? :developer : :coworker
  in_repository user do
    commands.split(';').each { |command| run command.strip }
  end
end


When(/^I run `([^`]+)` with the last( two)? commit shas?$/) do |command, two|
  count = two ? 2 : 1
  shas = recent_commit_shas(count).join(' ')
  step "I run `#{command} #{shas}`"
end


When(/^I run `(.+?)` and enter "(.+?)"$/) do |command, input|
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


Then(/^I get the error "(.+?)"$/) do |str|
  verify_error str
end


Then(/^I get the error$/) do |str|
  verify_error str.strip
end


Then(/^it runs no Git commands$/) do
  expect(commands_of_last_run).to be_empty
end


Then(/^it runs the Git commands$/) do |expected_steps|
  sha_regex = /\[SHA:(.+?)\]/

  # Replace SHA placeholders with the real SHAs
  expected_steps.map_column! 'COMMAND' do |command|
    command.gsub(sha_regex) do |sha_expression|
      commit_message = sha_expression.match(sha_regex).captures[0].strip
      output_of "git reflog --grep-reflog='commit: #{commit_message.strip}' --format='%H'"
    end
  end

  expected_steps.diff! commands_of_last_run.unshift(expected_steps.headers)
end


Then(/^I see no output$/) do
  expect(@last_run_result.out).to eql ''
end


Then(/^I don't see "(.*)"$/) do |string|
  expect(@last_run_result.out).not_to include(string)
end


Then(/^I see "(.*)"$/) do |string|
  actual = unformatted_last_run_output.strip
  expect(actual).to eql string
end


Then(/^I see$/) do |output|
  actual = unformatted_last_run_output.strip
  expect(actual).to eql output
end


Then(/^the output begins with "(.*)"$/) do |output|
  actual = unformatted_last_run_output
  expect(actual).to start_with output
end


Then(/^the output contains "(.*)"$/) do |output|
  actual = unformatted_last_run_output
  expect(actual).to include output
end


Then(/^I see the "(.+?)" man page$/) do |manpage|
  expect(@last_run_result.out).to eql "man called with: #{manpage}\n"
end
