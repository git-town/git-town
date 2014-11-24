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
  output = @last_run_result.out + @last_run_result.err
  expect(output).to include(error_message), %(
    "EXPECTED
    ***************************************************

    #{output.gsub '\n', "\n"}

    ***************************************************
    TO INCLUDE '#{error_message}'

  )
end
