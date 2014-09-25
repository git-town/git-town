require 'kappamaki'
require 'open4'
require 'rspec'

# The files to ignore when checking files
IGNORED_FILES = %w[ tags ]

Before do
  # Enter the test repository
  FileUtils.rm_r '/tmp/git_town_specs_remote', force: true
  Dir.mkdir '/tmp/git_town_specs_remote'
  run_this "git init --bare /tmp/git_town_specs_remote"

  FileUtils.rm_r '/tmp/git_town_specs', force: true
  Dir.chdir '/tmp'
  run_this 'git clone git_town_specs_remote git_town_specs'

  # Create the config file
  Dir.chdir '/tmp/git_town_specs'
  run_this 'git config git-town.main-branch-name main'

  # Create the master branch
  File.write '.gitignore', '.gittownrc'
  run_this "git add .gitignore ; git commit -m 'Initial commit' ; git push -u origin master"

  # Create the main branch
  run_this "git checkout -b main master ; git push -u origin main"

  def script_path(operation:, command:)
    "/tmp/#{operation.gsub ' ', '_'}_#{command}_#{Dir.pwd.gsub '/', '_'}"
  end
end


at_exit do
  FileUtils.rm_r '/tmp/git_town_specs_remote', force: true
  FileUtils.rm_r '/tmp/git_town_specs', force: true
end
