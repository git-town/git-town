When /^(?:Charlie|my coworker) runs `([^`]+)`$/ do |command|
  at_path coworker_repository_path do
    run command
  end
end


When /^I run `([^`]+)`( while allowing errors)?$/ do |command, allow_failures|
  run command, allow_failures: allow_failures
end


When /^I run `git extract refactor` with the last( two)? commit shas?( while allowing errors)?$/ do |two, allow_failures|
  count = two ? 2 : 1
  shas = recent_commit_shas(count).join(' ')
  run "git extract refactor #{shas}", allow_failures: allow_failures
end


When /^I run `(.+?)` and enter "(.*?)"$/ do |command, user_input|
  @result = run command, input: user_input, allow_failures: true
end


When /^I run `(.+?)` and then clear and exit the editor( while allowing errors)?$/ do |command, allow_failures|
  @result = run command, input: "dGZZ", allow_failures: allow_failures
end




Then /^I get the error "(.+?)"$/ do |error_message|
  expect(@last_run_result.status).to_not eq 0
  output = @last_run_result.out + @last_run_result.err
  expect(output).to include(error_message),
                    "EXPECTED\n\n***************************************************\n\n#{output.gsub '\n', "\n"}\n\n***************************************************\n\nTO INCLUDE '#{error_message}'\n"
end
