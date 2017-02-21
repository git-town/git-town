# frozen_string_literal: true
Then(/^I see a new (.+?) pull request for the "([^"]+)" branch in the "(.+?)" repo in my browser$/) do |domain, branch, repo|
  url = pull_request_url domain: domain, branch: branch, repo: repo
  expect(@last_run_result.out.strip).to include "#{@tool} called with: #{url}"
end


Then(/^I see a new (.+?) pull request for the "([^"]+)" branch against the "(.+?)" branch in the "(.+?)" repo in my browser$/) do |domain, child_branch, parent_branch, repo|
  url = pull_request_url domain: domain, branch: child_branch, parent_branch: parent_branch, repo: repo
  expect(@last_run_result.out.strip).to include "#{@tool} called with: #{url}"
end


Then(/^I see the (Bitbucket|GitHub|GitLab) homepage of the "(.+?)" repository in my browser$/) do |domain, repository|
  expect(@last_run_result.out.strip).to eql "#{@tool} called with: #{repository_homepage_url domain, repository}"
end
