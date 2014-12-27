When(/^(?:Charlie|my coworker) runs `([^`]+)`$/) do |command|
  at_path coworker_repository_path do
    run command
  end
end


When(/^I run `([^`]+)`( while allowing errors)?$/) do |commands, allow_failures|
  commands.split(';').each do |command|
    run command.strip, allow_failures: allow_failures
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


Then(/^I see a new (.+?) pull request for the "(.+?)" branch in my browser$/) do |domain, branch_name|
  expect(@last_run_result.out).to eql "#{@tool} called with: #{pull_request_url domain, branch_name}\n"
end


Then(/^I see the homepage of my (.+?) repository in my browser$/) do |domain|
  expect(@last_run_result.out).to eql "#{@tool} called with: #{repository_homepage_url domain}\n"
end


Then(/^it runs no Git commands$/) do
  expect(commands_of_last_run).to be_empty
end


Then(/^it runs the Git commands$/) do |expected_steps|

  # Replace SHA placeholders with the real SHAs of the respective branches
  expected_steps.map_column! 'COMMAND' do |command|
    command.gsub(/\[\[.+?\]\]/) do |sha_expression|
      branch_name = sha_expression.match(/"(.+?)" branch SHA/).captures[0]
      fail 'No branch name found' unless branch_name
      sha_of_branch branch_name
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
