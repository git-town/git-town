Given(/^I already have the Git autocompletion symlink$/) do
  FileUtils.rm FISH_AUTOCOMPLETIONS_PATH
  FileUtils.symlink 'foo',
                    FISH_AUTOCOMPLETIONS_PATH
end


Given(/^I have an existing Git autocompletion file$/) do
  FileUtils.rm FISH_AUTOCOMPLETIONS_PATH
  IO.write FISH_AUTOCOMPLETIONS_PATH,
           'existing Git autocompletion data'
end


Given(/^I have no fish autocompletion folder$/) do
  FileUtils.rm_rf File.expand_path('~/.config/fish/completions')
end


Given(/^I have an empty fish autocompletion folder$/) do
  FileUtils.rm_r FISH_AUTOCOMPLETIONS_PATH
end




Then(/^I still have my original Git autocompletion file$/) do
  expect(IO.read FISH_AUTOCOMPLETIONS_PATH).to eql 'existing Git autocompletion data'
end


Then(/^I still have my original Git autocompletion symlink$/) do
  expect(File.symlink? FISH_AUTOCOMPLETIONS_PATH).to be_truthy
end
