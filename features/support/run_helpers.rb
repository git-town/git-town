def output_of command
  run(command).out.strip
end


def array_output_of command
  output_of(command).split("\n").map(&:strip)
end


def integer_output_of command
  output_of(command).to_i
end


def git_town_command? command
  %w(extract hack prune-branches ship sync-fork sync kill).any? do |subcommand|
    command.starts_with? "git #{subcommand}"
  end
end


def print_result result
  puts ''
  puts "#{result.location}$ #{result.command}"
  puts "#{result.out}#{result.err}"
  puts ''
end


def run command, allow_failures: false, debug: false, input: nil
  result = run_shell_command command, input
  should_error = result.error && !allow_failures

  print_result(result) if should_error || should_print_command_output?(command, debug)
  fail 'Command not successful!' if should_error

  @last_run_result = result
end


def run_shell_command command, input
  result = OpenStruct.new(command: command, location: Dir.pwd.split(/[_\/]/).last)

  status = Open4.popen4(command) do |_pid, stdin, stdout, stderr|
    stdin.puts input if input
    stdin.close
    result.out = stdout.read
    result.err = stderr.read
  end

  result.error = status.exitstatus != 0
  result
end



def should_print_command_output? command, debug
  debug || ENV['DEBUG'] || (ENV['DEBUG_COMMANDS'] && git_town_command?(command))
end
