require 'active_support/all'
require 'kappamaki'
require 'open4'
require 'rspec'

SOURCE_DIRECTORY = "#{File.dirname(__FILE__)}/../../src"
SHELL_OVERRIDE_DIRECTORY = "#{File.dirname(__FILE__)}/shell_overrides"

REPOSITORY_BASE = Dir.mktmpdir
TOOLS_INSTALLED_FILENAME = "#{REPOSITORY_BASE}/tools_installed.txt"

Before do
  Dir.chdir REPOSITORY_BASE
  FileUtils.rm_rf Dir.glob("#{REPOSITORY_BASE}/*")

  # Create origin repository
  create_repository :origin

  # Create the local repository
  clone_repository :origin, :developer

  # Create the main branch
  in_repository :developer do
    run 'touch .gitignore ; git add .gitignore ; git commit -m "Initial commit"; git push -u origin master'
    run 'git checkout -b main master ; git push -u origin main'
  end

  # Set main as the default branch
  in_repository :origin do
    run 'git symbolic-ref HEAD refs/heads/main'
  end

  # Fetch the default branch, delete master
  in_repository :developer do
    run 'git fetch'
    run 'git push origin :master'
    run 'git branch -d master'
  end

  go_to_repository :developer
end


After do
  expect(@result_with_unexpected_error).to be_nil
end

After '~@finishes-with-non-empty-stash' do
  expect(stash_size).to eql(0), 'Finished with non empty stash'
end


at_exit do
  FileUtils.rm_rf REPOSITORY_BASE
end
