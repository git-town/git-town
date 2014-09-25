require 'kappamaki'
require 'open4'
require 'rspec'

# The files to ignore when checking files
IGNORED_FILES = %w[ tags ]

Before do
  # Enter the test repository
  FileUtils.rm_r '/tmp/git_town_specs_remote', force: true
  Dir.mkdir '/tmp/git_town_specs_remote'
  run "git init --bare /tmp/git_town_specs_remote"

  FileUtils.rm_r '/tmp/git_town_specs', force: true
  Dir.chdir '/tmp'
  run 'git clone git_town_specs_remote git_town_specs'

  # Create the config file
  Dir.chdir '/tmp/git_town_specs'
  run 'git config git-town.main-branch-name main'

  # Create the master branch
  File.write '.gitignore', '.gittownrc'
  run "git add .gitignore ; git commit -m 'Initial commit' ; git push -u origin master"

  # Create the main branch
  run "git checkout -b main master ; git push -u origin main"
end


at_exit do
  FileUtils.rm_r '/tmp/git_town_specs_remote', force: true
  FileUtils.rm_r '/tmp/git_town_specs', force: true
end
