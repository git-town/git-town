require 'active_support/all'
require 'kappamaki'
require 'open4'
require 'rspec'

SHELL_OVERRIDE_DIRECTORY = "#{File.dirname(__FILE__)}/shell_overrides"

# The files to ignore when checking files
IGNORED_FILES = %w(tags)

REPOSITORY_BASE = Dir.mktmpdir

Before do
  Dir.chdir REPOSITORY_BASE
  FileUtils.rm_rf Dir.glob("#{REPOSITORY_BASE}/*")

  # Create remote repository
  create_repository remote_repository_path

  # Create the local repository
  clone_repository remote_repository_path, local_repository_path
  at_path local_repository_path do
    # Create the master branch
    run 'touch .gitignore ; git add .gitignore ; git commit -m "Initial commit"; git push -u origin master'

    # Create the main branch
    run 'git checkout -b main master ; git push -u origin main'

    # Configuration
    run 'git config git-town.main-branch-name main'
    run 'git config git-town.non-feature-branch-names ""'
    run 'git config push.default simple'
    run 'git config core.editor vim'
  end

  # Set the default branch
  at_path remote_repository_path do
    run 'git symbolic-ref HEAD refs/heads/main'
  end

  # Fetch the default branch
  at_path local_repository_path do
    run 'git fetch'
  end

  # Create the coworker repository
  clone_repository remote_repository_path, coworker_repository_path
  at_path coworker_repository_path do
    run 'git config git-town.main-branch-name main'
  end

  Dir.chdir local_repository_path
end


After '~@finishes-with-non-empty-stash' do
  expect(stash_size).to eql(0), 'Finished with non empty stash'
end


at_exit do
  FileUtils.rm_rf REPOSITORY_BASE
end
