When /^Charlie runs `([^`]+)`$/ do |command|
  at_path coworker_repository_path do
    @last_run_result = run command
  end
end


When /^I run `([^`]+)`$/ do |command|
  @last_run_result = run command
end


When /^I run `([^`]+)` while allowing errors$/ do |command|
  @last_run_result = run command, allow_failures: true
end


When /^I run `git extract refactor` with the last commit sha as an argument$/ do
  sha = run("git log -n 1 | grep '^commit' | cut -d ' ' -f 2")[:out]
  @last_run_result = run "git extract refactor #{sha}"
end


When /^I run `git extract refactor` with the last commit sha as an argument while allowing errors$/ do
  sha = run("git log -n 1 | grep '^commit' | cut -d ' ' -f 2")[:out]
  @last_run_result = run "git extract refactor #{sha}", allow_failures: true
end


When /^I run `(.+?)` and enter "(.*?)"$/ do |command, user_input|
  @result = run command, input: user_input, allow_failures: true
end




Then /^I don't see "(.*?)"$/ do |unexpected_output|
  expect(@last_run_result[:out]).to_not include unexpected_output
end


Then /^I get the error "(.+?)"$/ do |error_message|
  expect(@last_run_result[:status]).to_not eq 0
  output = @last_run_result[:out] + @last_run_result[:err]
  expect(output).to include(error_message),
                    "EXPECTED\n\n***************************************************\n\n#{output.gsub '\n', "\n"}\n\n***************************************************\n\nTO INCLUDE '#{error_message}'\n"
end


Then /^I see "(.+?)"$/ do |expected_output|
  expect(@last_run_result[:out]).to include expected_output
end


Then /^It doesn't run the command "(.*?)"$/ do |unexpected_command|
  expect(@last_run_result[:out]).to_not include "#{unexpected_command}\n"
end


Then /^show me the output$/  do
  puts @last_run_result[:out]
  puts @last_run_result[:err]
end
