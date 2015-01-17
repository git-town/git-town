require 'active_support/all'
require 'kappamaki'
require 'open4'
require 'rspec'

SOURCE_DIRECTORY = "#{File.dirname(__FILE__)}/../../src"
SHELL_OVERRIDE_DIRECTORY = "#{File.dirname(__FILE__)}/shell_overrides"

# MEMOIZED_REPOSITORY_BASE = Dir.mktmpdir 'memoized'
REPOSITORY_BASE = Dir.mktmpdir
TOOLS_INSTALLED_FILENAME = "#{REPOSITORY_BASE}/tools_installed.txt"

class String
  # colorization
  def colorize(color_code)
    "\e[#{color_code}m#{self}\e[0m"
  end

  def yellow
    colorize(43)
  end
end

class MemSingle
  def initialize
    @log = Dir.mktmpdir 'memoized2'
  end

  @@instance = MemSingle.new

  def self.instance
    return @@instance
  end

  def val
    puts "Requested MEM DIR by #{caller_locations(1,1)[0].label} by #{caller_locations(2,1)[0].label} by #{caller_locations(3,1)[0].label}".yellow
    @log
  end

  private_class_method :new
end



def setup_environment
  FileUtils.rm_rf Dir.glob("#{REPOSITORY_BASE}/*")
  FileUtils.cp_r "#{MemSingle.instance.val}/.", REPOSITORY_BASE
  p 'NO MORE MEM ACCESS AFTER THIS POINT'

  Dir.chdir REPOSITORY_BASE
  go_to_repository :developer
end

def memoize_environment
  Dir.chdir MemSingle.instance.val
  FileUtils.rm_rf Dir.glob("#{MemSingle.instance.val}/*")

  # Create origin repository
  create_repository :origin, memoized: true

  # Create the local repository (~1/3)
  clone_repository :origin, :developer, memoized: true

  # Set main as the default branch
  in_repository :origin, memoized: true do
    run 'git symbolic-ref HEAD refs/heads/main'
  end

  in_repository :developer, memoized: true do
    # Create the main branch (~1/3)
    run 'touch .gitignore ; git add .gitignore ; git commit -m "Initial commit"; git push -u origin master'
    run 'git checkout -b main master ; git push -u origin main'

    # Fetch the default branch, delete master (~1/3)
    run 'git fetch'
    run 'git push origin :master'
    run 'git branch -d master'
  end

  $memoization_complete = true
  p 'DONE MEMOIZING ENVIRONMENT'
end


# rubocop:disable Style/GlobalVars
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
  FileUtils.rm_rf MemSingle.instance.val
end

# start = Time.now
# finish = Time.now
# x = finish - start
# p "AA [#{x}]"
