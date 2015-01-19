require 'active_support/all'
require 'kappamaki'
require 'open4'
require 'rspec'

# rubocop:disable all

SOURCE_DIRECTORY = "#{File.dirname(__FILE__)}/../../src"
SHELL_OVERRIDE_DIRECTORY = "#{File.dirname(__FILE__)}/shell_overrides"

MEMOIZED_REPOSITORY_BASE = Dir.mktmpdir 'memoized'
REPOSITORY_BASE = Dir.mktmpdir
TOOLS_INSTALLED_FILENAME = "#{REPOSITORY_BASE}/tools_installed.txt"


# copy entire contents of MEMOIZED_REPOSITORY_BASE to REPOSITORY_BASE
def setup_environment
  FileUtils.rm_rf Dir.glob("#{REPOSITORY_BASE}/*")
  FileUtils.cp_r "#{MEMOIZED_REPOSITORY_BASE}/.", REPOSITORY_BASE

  Dir.chdir REPOSITORY_BASE
  go_to_repository :developer
end


def memoize_environment
  FileUtils.rm_rf Dir.glob("#{MEMOIZED_REPOSITORY_BASE}/*")
  FileUtils.rm_rf Dir.glob("#{REPOSITORY_BASE}/*")

  # Create origin repository
  create_repository :origin

  # Create the local repository
  clone_repository :origin, :developer

  # Set main as the default branch
  in_repository :origin do
    run 'git symbolic-ref HEAD refs/heads/main'
  end

  in_repository :developer do
    # Create the main branch
    run 'touch .gitignore ; git add .gitignore ; git commit -m "Initial commit"; git push -u origin master'
    run 'git checkout -b main master ; git push -u origin main'

    # Fetch the default branch, delete master
    run 'git fetch'
    run 'git push origin :master'
    run 'git branch -d master'
  end

  FileUtils.cp_r "#{REPOSITORY_BASE}/.", MEMOIZED_REPOSITORY_BASE

  $memoization_complete = true
end


Before do
  $memoization_complete ||= false
  memoize_environment unless $memoization_complete
  setup_environment
end


After '~@finishes-with-non-empty-stash' do
  expect(stash_size).to eql(0), 'Finished with non empty stash'
end


at_exit do
  FileUtils.rm_rf REPOSITORY_BASE
  FileUtils.rm_rf MEMOIZED_REPOSITORY_BASE
end
