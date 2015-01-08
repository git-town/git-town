When(/^(I|my coworker) runs? `([^`]+)`( while allowing errors)?$/) do |who, commands, allow_failures|
  path = (who == 'I') ? local_repository_path : coworker_repository_path
  at_path(path) do
    commands.split(';').each do |command|
      run command.strip, allow_failures: allow_failures
    end
  end
end


When(/^I run `([^`]+)` with the last( two)? commit shas?( while allowing errors)?$/) do |command, two, allow_failures|
  count = two ? 2 : 1
  shas = recent_commit_shas(count).join(' ')
  run "#{command} #{shas}", allow_failures: allow_failures
end


When(/^I run `(.+?)` and enter (.+?)$/) do |command, input|
  inputs = prepare_user_input input
  @result = run command, inputs: inputs, allow_failures: true
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


Then(/^I see "(.*)"$/) do |string|
  expect(@last_run_result.out).to include string
end


Then(/^I see the (.+?) man page$/) do |manpage|
  expect(@last_run_result.out).to eql "man called with: #{manpage}\n"
end
