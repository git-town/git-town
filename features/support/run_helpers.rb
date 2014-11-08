def is_git_town_command? command
  %w(extract hack prune-branches ship sync-fork sync).any? do |subcommand|
    command.starts_with?("git #{subcommand}")
  end
end


def run command, allow_failures: false, debug: false, input: nil
  result = {}
  status = Open4::popen4(command) do |pid, stdin, stdout, stderr|
    stdin.puts input if input
    stdin.close
    result[:out] = stdout.read.strip
    result[:err] = stderr.read.strip
  end
  result[:status] = status.exitstatus
  if result[:status] != 0 and !allow_failures
    raise "\nCommand '#{command}' not successful! \n\n************\nOUT: '#{result[:out]}', \n\n************\nERR: '#{result[:err]}'\n\n"
  end
  if should_print_command_output?(command, debug)
    puts "\nRUNNING: #{command}"
    puts "#{result[:out]}\n"
    puts "#{result[:err]}\n"
  end
  OpenStruct.new result
end


def should_print_command_output? command, debug
  debug or ENV["DEBUG"] or (ENV['DEBUG_COMMANDS'] and is_git_town_command?(command))
end
