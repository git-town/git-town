# frozen_string_literal: true
Given(/^I have an existing Git autocompletion symlink$/) do
  symlinked_file = File.join(REPOSITORY_BASE, '.config/completions/custom')
  FileUtils.mkdir_p File.dirname(FISH_AUTOCOMPLETIONS_PATH)
  FileUtils.mkdir_p File.dirname(symlinked_file)
  FileUtils.symlink symlinked_file, FISH_AUTOCOMPLETIONS_PATH
  IO.write symlinked_file, 'existing Git autocompletion data'
end


Given(/^I have an existing Git autocompletion file$/) do
  FileUtils.mkdir_p File.dirname(FISH_AUTOCOMPLETIONS_PATH)
  IO.write FISH_AUTOCOMPLETIONS_PATH, 'existing Git autocompletion data'
end


Given(/^I have no fish autocompletion folder$/) do
  # empty for readability
end


Given(/^I have an empty fish autocompletion folder$/) do
  FileUtils.mkdir_p File.dirname(FISH_AUTOCOMPLETIONS_PATH)
end




Then(/^I am in the project root folder$/) do
  expect(Dir.pwd).to eql git_root_folder
end


Then(/^I have a Git autocompletion file$/) do
  expect(IO.read FISH_AUTOCOMPLETIONS_PATH).to include 'complete --command git'
end


Then(/^I still have my original Git autocompletion file$/) do
  expect(IO.read FISH_AUTOCOMPLETIONS_PATH).to eql 'existing Git autocompletion data'
end
