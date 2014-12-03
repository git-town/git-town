require 'English'

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


When(/^I run `(.+?)` and enter "(.*?)"$/) do |command, user_input|
  @result = run command, input: user_input, allow_failures: true
end


When(/^I run `(.+?)` and enter an empty commit message?$/) do |command|
  @result = run command, input: 'dGZZ', allow_failures: true
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


Then(/^the output should contain '(.*?)'$/) do |string|
  # @shell_overrides_export
  if string =~ /<#(.+)=(.+)#>/
    override_placeholder, override_variable, override_value = $LAST_MATCH_INFO.to_a

    expect(@last_run_result.out).to include string.gsub(override_placeholder, '')
    expect(@shell_overrides_export[override_variable]).to eql override_value
  else
    expect(@last_run_result.out).to include string
  end
end
