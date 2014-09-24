Given /^I don't have a main branch name configured$/ do
  run_this 'git config --unset git-town.main-branch-name'
end


Given /^I have an old configuration file with (.*?)$/ do |data|
  run_this 'git config --unset git-town.main-branch-name'
  data = Kappamaki.attributes_from_sentence data
  IO.write '.main_branch_name', data[:'main branch']
end



Then /^I don't have an old configuration file anymore$/ do
  expect(File.exist? '.main_branch_name').to be_falsy
end


Then /^the main branch name is now configured as "(.*?)"$/ do |main_branch_name|
  expect(run_this('git config --get git-town.main-branch-name')[:out]).to eql main_branch_name
end

