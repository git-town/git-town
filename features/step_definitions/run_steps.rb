When(/^(?:Charlie|my coworker) runs `([^`]+)`$/) do |command|
  at_path coworker_repository_path do
    run command
  end
end


When(/^I run `([^`]+)`( while allowing errors)?$/) do |command, allow_failures|
  run command, allow_failures: allow_failures
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


Then(/^I see a browser window for a new pull request on (.+) for the "(.+)" branch$/) do |domain, branch_name|
  expect(@last_run_result.out).to eql "open called with: #{remote_pull_request_url domain, branch_name}\n"
end


Then(/^I see "(.*)"$/) do |string|
  expect(@last_run_result.out).to include string
end
