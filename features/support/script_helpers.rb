# Returns the path to the abort or continue script for the given command
def script_path(operation:, command:)
  "/tmp/#{operation.gsub(' ', '_').gsub('-', '_')}_#{command}_#{Dir.pwd.gsub '/', '_'}"
end
