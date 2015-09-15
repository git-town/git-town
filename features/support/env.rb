require 'active_support/all'
require 'kappamaki'
require 'mortadella'
require 'open4'
require 'pathname'
require 'rspec'


SOURCE_DIRECTORY = "#{File.dirname(__FILE__)}/../../src"
GIT_TOWN_DIRECTORY = File.expand_path('..', SOURCE_DIRECTORY)
SHELL_OVERRIDE_DIRECTORY = "#{File.dirname(__FILE__)}/shell_overrides"

MEMOIZED_REPOSITORY_BASE = Dir.mktmpdir 'memoized'
REPOSITORY_BASE = Dir.mktmpdir
TOOLS_INSTALLED_FILENAME = "#{REPOSITORY_BASE}/tools_installed.txt"

FISH_AUTOCOMPLETIONS_PATH = File.expand_path '~/.config/fish/completions/git.fish'

DEBUG = {

  # Prints debug info for all activities
  all: ENV['DEBUG'],

  # Prints debug info only for the Git commands run
  commands_only: ENV['DEBUG_COMMANDS']
}


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
  @non_empty_stash_expected = false
end


After do
  unless @non_empty_stash_expected
    expect(stash_size).to eql(0), 'Finished with non empty stash'
  end
end


After '~@ignore-run-error' do
  if @last_run_result && !@error_expected
    puts unformatted_last_run_output if @last_run_result.error
    expect(@last_run_result.error).to be_falsy, 'Expected no runtime error'
  end
end


Before '@debug' do
  DEBUG[:all] = true
end

After '@debug' do
  DEBUG[:all] = ENV['DEBUG']
end


Before '@debug-commands' do
  DEBUG[:commands_only] = true
end

After '@debug-commands' do
  DEBUG[:commands_only] = ENV['DEBUG_COMMANDS']
end


Around '@modifies-fish-autocompletions' do |_scenario, block|
  completions_path = File.expand_path('~/.config/fish/completions')
  backup_path = File.expand_path('~/__config_fish_backup__')

  FileUtils.cp_r completions_path, backup_path

  block.call

  FileUtils.rm_rf completions_path
  FileUtils.cp_r backup_path, completions_path
  FileUtils.rm_rf backup_path
end


at_exit do
  FileUtils.rm_rf REPOSITORY_BASE
  FileUtils.rm_rf MEMOIZED_REPOSITORY_BASE
end
