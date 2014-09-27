require 'kappamaki'
require 'open4'
require 'rspec'

# The files to ignore when checking files
IGNORED_FILES = %w[ tags ]

Before do
  Dir.chdir repositiory_base_path

  create_repository remote_repository_path

  clone_repository remote_repository_path, local_repository_path do
    File.write '.gitignore', ''
    run 'git add .gitignore ; git commit -m "Initial commit"; git push -u origin master'
    run 'git checkout -b main master ; git push -u origin main'
    run 'git config git-town.main-branch-name main'
  end

  clone_repository remote_repository_path, coworker_repository_path do
    run 'git checkout main'
  end

  Dir.chdir local_repository_path
end


at_exit do
  Dir.chdir repositiory_base_path
  delete_repository remote_repository_path
  delete_repository local_repository_path
  delete_repository coworker_repository_path
end
