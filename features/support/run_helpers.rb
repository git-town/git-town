def array_output_of command
  output_of(command).split("\n").map(&:strip)
end


# Returns an array of the commands that were run in the last invocation of "run"
def commands_of_last_run_outside_git
  command_regex = /
    \[1m          # bold text
    (.+?)         # the command
    \s*           # any extra whitespace
    \n            # newline at the end
  /x
  @last_run_result.out.scan command_regex
end


# This regex parses Git commands out of console output.
# It is used in commands_of_last_run
GIT_COMMAND_REGEX = /
  \[1m          # bold text
  \[(.*?)\]     # branch name in square brackets
  \s            # space between branch name and Git command
  (.+?)         # the Git command
  \s*           # any extra whitespace
  \n            # newline at the end
/x


# Returns a Mortadella instance containing the Git commands run
# in the last invocation of "run".
def commands_of_last_run
  result = Mortadella.new headers: %w(BRANCH COMMAND), dry: 'BRANCH'
  @last_run_result.out.scan(GIT_COMMAND_REGEX).each { |command| result << command }
  result
end


def integer_output_of command
  output_of(command).to_i
end


def git_town_command? command
  %w(extract hack kill new-pull-request prune-branches rename-branch repo ship sync-fork sync town).any? do |subcommand|
    command.starts_with? "git #{subcommand}"
  end
end


def output_of command
  run(command).out.strip
end


def print_result result
  puts ''
  puts %(#{result.location}$ #{result.command}
         #{result.out}).gsub(/^/, '> ')
  puts ''
end


def run command, inputs: []
  result = run_shell_command command, inputs
  is_git_town_command = git_town_command? command
  raise_error = (!is_git_town_command && result.error) || result_has_shell_error?(result)

  print_result(result) if raise_error || should_print_command_output?(command)
  fail 'Command not successful!' if raise_error

  @last_run_result = result if is_git_town_command
  result
end


def result_has_shell_error? result
  # Shell errors have the format
  #   <filename>: line <line number>: <error message>
  result.out.include? File.expand_path('../../../src/', __FILE__)
end


def run_shell_command command, inputs = []
  result = OpenStruct.new(command: command, location: Pathname.new(Dir.pwd).basename)
  command = "#{shell_overrides}; #{command} 2>&1"
  kill = inputs.pop if inputs.last == '^C' # command shouldn't error if user aborts it

  status = Open4.popen4(command) do |_pid, stdin, stdout, _stderr|
    inputs.each { |input| stdin.puts input }
    stdin.close
    result.out = stdout.read
  end

  result.error = status.exitstatus != 0 && !kill
  result
end


def shell_overrides
  "PATH=#{SOURCE_DIRECTORY}:#{SHELL_OVERRIDE_DIRECTORY}:$PATH; export WHICH_SOURCE=#{TOOLS_INSTALLED_FILENAME}"
end


def should_print_command_output? command
  DEBUG[:all] || (DEBUG[:commands_only] && git_town_command?(command))
end


# Output of last `run` without text formatting
def unformatted_last_run_output
  @last_run_result.out
    .gsub(/\e[^m]*m/, '')   # remove color codes
    .gsub("\x0F", '')       # remove artifacts created by CircleCI
    .gsub(/\\u\d*F/, '')    # remove artifacts created by CircleCI
end
