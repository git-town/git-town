def run_this command, allow_failures: false, debug: false, input: nil
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
  if debug or $debug
    puts "\nRUNNING: #{command}"
    puts "#{result[:out]}\n"
    puts "#{result[:err]}\n"
  end
  result
end

