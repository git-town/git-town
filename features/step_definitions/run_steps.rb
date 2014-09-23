When /^I run `([^`]+)`$/ do |command|
  @last_run_result = run_this command
end


When /^I run `([^`]+)` while allowing errors$/ do |command|
  @last_run_result = run_this command, allow_failures: true
end


When /^I run `git extract refactor` with the last commit sha as an argument$/ do
  sha = run_this("git log -n 1 | grep '^commit' | cut -d ' ' -f 2")[:out]
  @last_run_result = run_this "git extract refactor #{sha}"
end


When /^I run `git extract refactor` with the last commit sha as an argument while allowing errors$/ do
  sha = run_this("git log -n 1 | grep '^commit' | cut -d ' ' -f 2")[:out]
  @last_run_result = run_this "git extract refactor #{sha}", allow_failures: true
end


When /^I run `(.+?)` and enter "(.*?)"$/ do |command, user_input|
  @result = run_this command, input: user_input, allow_failures: true
end



Then /^I get the error "(.+?)"$/ do |error_message|
  expect(@last_run_result[:status]).to_not eq 0
  output = @last_run_result[:out] + @last_run_result[:err]
  expect(output).to include error_message
end


Then /^show me the output$/  do
  puts @last_run_result[:out]
end
