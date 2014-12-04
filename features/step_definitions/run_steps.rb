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


Then(/^it runs no Git commands$/) do
  expect(@last_run_result.out.scan(/\[1m\[(.*?)\] (.*?)\n/)).to be_empty
end


Then(/^it runs the Git commands$/) do |steps_table|
  actual_steps = [['BRANCH', 'COMMAND']]
  actual_steps.concat @last_run_result.out.scan(/\[1m\[(.*?)\] (.*?)\n/)
  expected_steps = steps_table.raw
  expect(expected_steps.size).to eq(actual_steps.size),
                                 "expected #{expected_steps.size} steps, found #{actual_steps.size}: #{actual_steps}"
  actual_steps.each_with_index do |actual_step, i|
    expected_step = expected_steps[i]
    expect(actual_step[0]).to eq expected_step[0]
    expect(actual_step[1]).to match Regexp.new "^#{expected_step[1]}$"
  end
end
