Then(/^I see a new (.+?) pull request for the "([^"]+)" branch in the "(.+?)" repo in my browser$/) do |domain, branch, repo|
  expect(@last_run_result.out.strip).to eql "#{@tool} called with: #{pull_request_url domain, branch, repo}"
end


Then(/^I see a new (.+?) pull request for the "([^"]+)" branch against the "(.+?)" branch in the "(.+?)" repo in my browser$/) do |domain, child_branch, parent_branch, repo|
  expect(@last_run_result.out.strip).to eql "#{@tool} called with: #{pull_request_url domain, branch, repo}"
end


Then(/^I see the (Bitbucket|GitHub) homepage of the "(.+?)" repository in my browser$/) do |domain, repository|
  expect(@last_run_result.out.strip).to eql "#{@tool} called with: #{repository_homepage_url domain, repository}"
end
