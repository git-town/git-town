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
  if debug or ENV["DEBUG"]
    puts "\nRUNNING: #{command}"
    puts "#{result[:out]}\n"
    puts "#{result[:err]}\n"
  end
  OpenStruct.new result
end

def run_and_store command, allow_failures: false, debug: ENV['OUTPUT_COMMANDS'], input: nil
  @last_run_result = run command, allow_failures: allow_failures, debug: debug, input: input
end

