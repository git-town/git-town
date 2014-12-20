def array_output_of command
  output_of(command).split("\n").map(&:strip)
end


# Returns an array of the Git commands that were run in the last invocation of "run"
# with the form [<branch_name>, <command>]
def commands_of_last_run
  command_regex = /
    \[1m          # bold text
    \[(.*?)\]     # branch name in square brackets
    \s            # space between branch name and Git command
    (.*?)         # the Git command
    \n            # newline at the end
  /x
  @last_run_result.out.scan command_regex
end


def integer_output_of command
  output_of(command).to_i
end


def git_town_command? command
  %w(extract hack kill pr prune-branches repo ship sync-fork sync town).any? do |subcommand|
    command.starts_with? "git #{subcommand}"
  end
end


def output_of command
  run(command).out.strip
end


def prepare_user_input input
  if input == 'an empty commit message'
    ['dGZZ']
  else
    Kappamaki.from_sentence(input)
  end
end


def print_result result
  puts ''
  puts "#{result.location}$ #{result.command}"
  puts "#{result.out}"
  puts ''
end


def run command, allow_failures: false, debug: false, inputs: []
  result = run_shell_command command, inputs
  should_error = should_error? result, allow_failures

  print_result(result) if should_error || should_print_command_output?(command, debug)
  fail 'Command not successful!' if should_error

  @last_run_result = result if git_town_command?(command)

  result
end


def result_has_shell_error? result
  # Shell errors have the format
  #   <filename>: line <line number>: <error message>
  result.out.include? File.expand_path('../../../src/', __FILE__)
end


def run_shell_command command, inputs
  result = OpenStruct.new(command: command, location: Dir.pwd.split(/[_\/]/).last)
  command = "PATH=#{SHELL_OVERRIDE_DIRECTORY}:$PATH; #{command} 2>&1"

  status = Open4.popen4(command) do |_pid, stdin, stdout, _stderr|
    inputs.each { |input| stdin.puts input }
    stdin.close
    result.out = stdout.read
  end

  result.error = status.exitstatus != 0
  result
end


def should_error? result, allow_failures
  (result.error && !allow_failures) || result_has_shell_error?(result)
end


def should_print_command_output? command, debug
  debug || ENV['DEBUG'] || (ENV['DEBUG_COMMANDS'] && git_town_command?(command))
end
