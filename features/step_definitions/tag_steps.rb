# frozen_string_literal: true

Given(/^I have the following tags$/) do |tags|
  tags.hashes.each do |tag|
    send "create_#{tag['LOCATION']}_tag", tag['NAME']
  end
end


Given(/^I have a remote tag "([^"]+)" that is not on a branch$/) do |name|
  create_standalone_remote_tag name
end




Then(/^I now have the following tags$/) do |expected_tags|
  expected_tags.diff! TagFinder.all_tags.to_table
end
