require 'kappamaki'
require 'open4'
require 'rspec'

# The files to ignore when checking files
IGNORED_FILES = %w[ tags ]

Before do
  Dir.chdir repositiory_base

  create_repository remote_repository

  clone_repository remote_repository, local_repository

  in_repository local_repository do
    File.write '.gitignore', ''
    run 'git add .gitignore ; git commit -m "Initial commit"; git push -u origin master'
    run 'git checkout -b main master ; git push -u origin main'
    run 'git config git-town.main-branch-name main'
  end

  clone_repository remote_repository, coworker_repository

  in_repository coworker_repository do
    run 'git checkout main'
  end

  Dir.chdir local_repository
end


at_exit do
  Dir.chdir repositiory_base
  delete_repository remote_repository
  delete_repository local_repository
  delete_repository coworker_repository
end
