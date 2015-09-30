def array_output_of command, ignore_errors: false
  output_of(command, ignore_errors: ignore_errors).split("\n").map(&:strip)
end


COMMAND_REGEX = /
  ^
  \e\[1m                     # bold text
  (?:\[(.+?)\]\s)?           # branch name in square brackets
  ([[:graph:]][[:print:]]+?) # the command - no leading whitespace
  \s*                        # trailing whitespace
  $
/x


# Returns an array of the commands that were run in the last invocation of "run"
def commands_of_last_run with_branch: true
  options = with_branch ? { headers: %w(BRANCH COMMAND), dry: 'BRANCH' } : { headers: %w(COMMAND) }
  result = Mortadella.new options
  @last_run_result.out.split("\n").each do |line|
    match = line.match COMMAND_REGEX
    next unless match
    row = [match[2]]
    row.unshift(match[1] || '<none>') if with_branch
    result << row
  end
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


def output_of command, ignore_errors: false
  run(command, ignore_errors: ignore_errors).out.strip
end


def print_result result
  puts ''
  puts %(#{result.location}$ #{result.command}
         #{result.out}).gsub(/^/, '> ')
  puts ''
end


def run command, inputs: [], ignore_errors: false
  result = run_shell_command command, inputs
  is_git_town_command = git_town_command? command
  raise_error = should_raise_error? is_git_town_command: is_git_town_command,
                                    result: result,
                                    ignore_errors: ignore_errors
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
  "PATH=#{SOURCE_DIRECTORY}:#{SHELL_OVERRIDE_DIRECTORY}:$PATH;"\
  "export WHICH_SOURCE=#{TOOLS_INSTALLED_FILENAME};"\
  "export GT_ENV=test"
end


def should_print_command_output? command
  DEBUG[:all] || (DEBUG[:commands_only] && git_town_command?(command))
end


# Returns whether a test should raise an error in the given situation
def should_raise_error? is_git_town_command:, result:, ignore_errors:
  ((!is_git_town_command && result.error) || result_has_shell_error?(result)) && !ignore_errors
end


# Output of last `run` without text formatting
def unformatted_last_run_output
  @last_run_result.out
    .gsub(/\e[^m]*m/, '')   # remove color codes
    .delete("\x0F")         # remove artifacts created by CircleCI
    .gsub(/\\u\d*F/, '')    # remove artifacts created by CircleCI
end
