require 'active_support/all'
require 'kappamaki'
require 'open4'
require 'rspec'


SOURCE_DIRECTORY = "#{File.dirname(__FILE__)}/../../src"
SHELL_OVERRIDE_DIRECTORY = "#{File.dirname(__FILE__)}/shell_overrides"

MEMOIZED_REPOSITORY_BASE = Dir.mktmpdir 'memoized'
REPOSITORY_BASE = Dir.mktmpdir
TOOLS_INSTALLED_FILENAME = "#{REPOSITORY_BASE}/tools_installed.txt"


# load memoized environment by copying contents
# of MEMOIZED_REPOSITORY_BASE to REPOSITORY_BASE
def setup_environment
  FileUtils.rm_rf REPOSITORY_BASE
  FileUtils.cp_r "#{MEMOIZED_REPOSITORY_BASE}/.", REPOSITORY_BASE
end


def initialize_environment
  # Create origin repo and set "main" as default branch
  create_repository :origin do
    run 'git symbolic-ref HEAD refs/heads/main'
  end

  clone_repository :origin, :developer

  # Initialize main branch
  in_repository :developer do
    run 'git checkout --orphan main'
    run 'git commit --allow-empty -m "Initial commit"'
    run 'git push -u origin main'
  end

  # memoize environment by saving directory contents
  FileUtils.cp_r "#{REPOSITORY_BASE}/.", MEMOIZED_REPOSITORY_BASE
end


AfterConfiguration do
  initialize_environment
end


Before do
  setup_environment
  go_to_repository :developer
end


Before do
  @error_expected = false
end


After do
  expect(@last_run_result.try :error).to be_falsy, 'Expected no runtime error' unless @error_expected
end


After '~@finishes-with-non-empty-stash' do
  expect(stash_size).to eql(0), 'Finished with non empty stash'
end


at_exit do
  FileUtils.rm_rf REPOSITORY_BASE
  FileUtils.rm_rf MEMOIZED_REPOSITORY_BASE
end
