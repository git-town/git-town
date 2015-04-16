Given(/^I have an uncommitted file$/) do
  @uncommitted_file_name = create_uncommitted_file
end


Given(/^I have an uncommitted file with name: "(.+?)" and content: "(.+?)"$/) do |name, content|
  create_uncommitted_file name: name, content: content
end


Given(/^I resolve the conflict in "(.+?)"$/) do |file_name|
  IO.write file_name, 'resolved content'
  run "git add #{file_name}"
end




Then(/^I (?:still|again) have my uncommitted file$/) do
  verify_uncommitted_file name: @uncommitted_file_name
end


Then(/^(?:now I|I still) have the following committed files$/) do |files_data|
  files_data.diff! files_in_branches
end


Then(/^I don't have any uncommitted files$/) do
  expect(uncommitted_files).to be_empty
end


Then(/^my workspace (?:still|again) has an uncommitted file with name: "([^"]+)" and content: "([^"]+)"$/) \
    do |name, content|
  verify_uncommitted_file name: name, content: content
end


Then(/^my uncommitted file is stashed$/) do
  expect(uncommitted_files).to_not include @uncommitted_file_name
  expect(stash_size).to eql 1
  @finishes_with_non_empty_stash = true
end
