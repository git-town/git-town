Given(/^I already have the Git autocompletion symlink$/) do
  FileUtils.rm File.expand_path('~/.config/fish/completions/git.fish')
  FileUtils.symlink 'foo',
                    File.expand_path('~/.config/fish/completions/git.fish')
end


Given(/^I have an existing Git autocompletion file$/) do
  FileUtils.rm File.expand_path('~/.config/fish/completions/git.fish')
  IO.write File.expand_path('~/.config/fish/completions/git.fish'),
           'existing Git autocompletion data'
end


Given(/^I have no fish autocompletion folder$/) do
  FileUtils.rm_rf File.expand_path('~/.config/fish/completions')
end


Given(/^I have an empty fish autocompletion folder$/) do
  suppress do
    FileUtils.rm File.expand_path('~/.config/fish/completions/git.fish')
  end
end




  expect(IO.read File.expand_path('~/.config/fish/completions/git.fish'))
    .to eql 'existing Git autocompletion data'
Then(/^I still have my original Git autocompletion file$/) do
end


  expect(File.symlink? File.expand_path('~/.config/fish/completions/git.fish')).to be_truthy
Then(/^I still have my original Git autocompletion symlink$/) do
end
