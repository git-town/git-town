# frozen_string_literal: true

require 'active_support/all'
require 'kappamaki'
require 'mortadella'
require 'open4'
require 'pathname'
require 'rspec'
require 'tmpdir'


SOURCE_DIRECTORY = File.join(File.dirname(__FILE__), '..', '..', 'src')
GIT_TOWN_DIRECTORY = File.expand_path('..', SOURCE_DIRECTORY)
SHELL_OVERRIDE_DIRECTORY = File.join(File.dirname(__FILE__), 'shell_overrides')

MEMOIZED_REPOSITORY_BASE = Dir.mktmpdir 'memoized'
REPOSITORY_BASE = Dir.mktmpdir
TOOLS_INSTALLED_FILENAME = File.join(REPOSITORY_BASE, 'tools_installed.txt')

FISH_AUTOCOMPLETIONS_PATH = File.join(REPOSITORY_BASE, '.config/fish/completions/git.fish')


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
  @debug = ENV['DEBUG']
  @debug_commands = ENV['DEBUG_COMMANDS']
  @temporary_shell_overrides_directory = Dir.mktmpdir 'temp_shell_overrides'
end


Before '@debug' do
  @debug = true
end


Before '@debug-commands' do
  @debug_commands = true
end


After do
  run 'git-town discard'
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


After do
  FileUtils.rm_rf @temporary_shell_overrides_directory
end


at_exit do
  FileUtils.rm_rf REPOSITORY_BASE
  FileUtils.rm_rf MEMOIZED_REPOSITORY_BASE
end
