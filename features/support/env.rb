require 'kappamaki'
require 'open4'
require 'rspec'

# The files to ignore when checking files
IGNORED_FILES = %w[ tags ]

Before do
  Dir.chdir repositiory_base

  # Create the remote repository
  create_repository remote_repository_path

  # Create the local repository
  clone_repository remote_repository_path, local_repository_path
  at_path local_repository_path do

    # Create the master branch
    run 'touch .gitignore ; git add .gitignore ; git commit -m "Initial commit"; git push -u origin master'

    # Create the main branch
    run 'git checkout -b main master ; git push -u origin main'
    run 'git config git-town.main-branch-name main'
  end

  # Create the coworker repository
  clone_repository remote_repository_path, coworker_repository_path
  at_path coworker_repository_path do
    run 'git checkout main'
  end

  Dir.chdir local_repository_path
end


at_exit do
  Dir.chdir repositiory_base
  delete_repository remote_repository_path
  delete_repository local_repository_path
  delete_repository coworker_repository_path
end
