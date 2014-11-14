def output_of command
  run(command).out.strip
end


def array_output_of command
  output_of(command).split("\n").map(&:strip)
end


def integer_output_of command
  output_of(command).to_i
end


def is_git_town_command? command
  %w(extract hack prune-branches ship sync-fork sync kill).any? do |subcommand|
    command.starts_with? "git #{subcommand}"
  end
end


def run command, allow_failures: false, debug: false, input: nil
  result = {}
  status = Open4::popen4(command) do |pid, stdin, stdout, stderr|
    stdin.puts input if input
    stdin.close
    result[:out] = stdout.read
    result[:err] = stderr.read
  end
  result[:status] = status.exitstatus
  if result[:status] != 0 and !allow_failures
    raise "\nCommand '#{command}' not successful! \n\n************\nOUT:\n#{result[:out]} \n\n************\nERR:\n#{result[:err]}\n\n"
  end
  if should_print_command_output? command, debug
    puts "\nRUNNING: #{command}"
    puts "#{result[:out]}\n"
    puts "#{result[:err]}\n"
  end
  @last_run_result = OpenStruct.new result
end


def should_print_command_output? command, debug
  debug or ENV["DEBUG"] or (ENV['DEBUG_COMMANDS'] and is_git_town_command? command)
end
