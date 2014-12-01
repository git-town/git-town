# Returns whether the given abort/continue script exists
def script_exists? *args
  File.exist? script_path(*args)
end


# Returns the path to the abort or continue script for the given command
def script_path command, action
  action = 'command_list' if action == 'continue'
  "/tmp/#{command.gsub(' ', '-')}_#{action}_#{Dir.pwd.gsub '/', '_'}"
end
