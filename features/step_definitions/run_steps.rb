When(/^(I|my coworker) runs? `([^`]+)`( it errors)?$/) do |who, commands, should_error|
  user = (who == 'I') ? :developer : :coworker
  in_repository user do
    commands.split(';').each do |command|
      run command.strip, should_error: should_error
      expect(@last_run_result.error).to be_truthy if should_error
    end
  end
end


When(/^I run `([^`]+)` with the last( two)? commit shas?( it errors)?$/) do |command, two, should_error|
  count = two ? 2 : 1
  shas = recent_commit_shas(count).join(' ')
  step "I run `#{command} #{shas}`#{should_error}"
end


When(/^I run `(.+?)` and enter "(.+?)"( it errors)?$/) do |command, input, should_error|
  @result = run command, input: input, should_error: should_error
  expect(@last_run_result.error).to be_truthy if should_error
end


When(/^I run `(.+?)` and enter an empty commit message it errors$/) do |command|
  step "I run `#{command}` and enter \"dGZZ\" it errors"
end




Then(/^I get the error "(.+?)"$/) do |error_message|
  expect(@last_run_result.error).to be_truthy
  expect(@last_run_result.out).to include(error_message), %(
    "EXPECTED
    ***************************************************

    #{@last_run_result.out.gsub '\n', "\n"}

    ***************************************************
    TO INCLUDE '#{error_message}'

  )
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
  actual = unformatted_last_run_output
  expect(actual).to eql "#{output}\n"
end


Then(/^I see the "(.+?)" man page$/) do |manpage|
  expect(@last_run_result.out).to eql "man called with: #{manpage}\n"
end
