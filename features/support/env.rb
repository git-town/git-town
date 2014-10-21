require 'kappamaki'
require 'open4'
require 'rspec'
require 'active_support/all'

# The files to ignore when checking files
IGNORED_FILES = %w[ tags ]

Before do
  Dir.chdir repositiory_base

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
  end

  # Set the default branch
  at_path remote_repository_path do
    run 'git symbolic-ref HEAD refs/head/main'
  end

  # Pull the default branch
  at_path local_repository_path do
    run 'git pull'
  end

  # Create the coworker repository
  clone_repository remote_repository_path, coworker_repository_path
  at_path coworker_repository_path do
    run 'git checkout main'
    run 'git config git-town.main-branch-name main'
  end

  Dir.chdir local_repository_path
end


Before '@github_query' do
  github_check_rate_limit!
end


After '~@finishes-with-non-empty-stash' do
  expect(stash_size).to eql 0
end


at_exit do
  Dir.chdir repositiory_base
  delete_repository remote_repository_path
  delete_repository local_repository_path
  delete_repository coworker_repository_path
end
