Given(/^my repo ignores files named "([^"]*)"$/) do |filename|
  create_local_commit branch: current_branch_name,
                      file_name: '.gitignore',
                      file_content: filename,
                      message: 'ignoring files'
end


Given(/^I have an uncommitted file(?: with name: "(.+?)" and content: "(.+?)")?$/) do |name, content|
  @uncommitted_file_name = name || 'uncommitted_file'
  @uncommitted_file_content = content || 'uncommitted content'
  create_uncommitted_file name: @uncommitted_file_name, content: @uncommitted_file_content
end


Given(/^I resolve the conflict in "([^"]+)"(?: with "([^"]+)")?$/) do |file_name, file_content|
  IO.write file_name, (file_content || 'resolved content')
  run "git add #{file_name}"
end




Then(/^I (?:still|again) have my uncommitted file$/) do
  verify_uncommitted_file name: @uncommitted_file_name, content: @uncommitted_file_content
end


Then(/^(?:now I|I still) have the following committed files$/) do |files_data|
  files_data.diff! files_in_branches
end


Then(/^I don't have any uncommitted files$/) do
  expect(uncommitted_files).to be_empty
end


Then(/^my uncommitted file is stashed$/) do
  expect(uncommitted_files).to_not include @uncommitted_file_name
  expect(stash_size).to eql 1
  @non_empty_stash_expected = true
end


Then(/^my workspace still contains the file "([^"]*)" with content "([^"]*)"$/) do |filename, content|
  verify_file filename, content
end
