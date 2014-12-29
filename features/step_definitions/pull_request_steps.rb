Then(/^I see a new (.+?) pull request for the "(.+?)" branch in the "(.+?)" repo in my browser$/) do |domain,
                                                                                                      branch,
                                                                                                      repo|
  expect(@last_run_result.out).to eql "#{@tool} called with: #{pull_request_url domain, branch, repo}\n"
end
