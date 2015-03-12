Given(/^I have an uncommitted file with name: "(.+?)" and content: "(.+?)"$/) do |name, content|
  IO.write name, content
end


Given(/^I resolve the conflict in "(.+?)"$/) do |file_name|
  IO.write file_name, 'resolved content'
  run "git add #{file_name}"
end




Then(/^(?:now I|I still) have the following committed files$/) do |files_data|
  files_data.diff! files_in_branches_array
end


Then(/^I don't have an uncommitted file with name: "(.+?)"$/) do |file_name|
  expect(uncommitted_files).to_not include file_name
end


Then(/^I don't have any uncommitted files$/) do
  expect(uncommitted_files).to be_empty
end


Then(/^I (?:still|again) have an uncommitted file with name: "([^"]+)" and content: "([^"]+)"$/) do |file_name, content|
  expect(uncommitted_files).to eql [file_name]
  expect(IO.read file_name).to eql content
end
