# Returns whether the given abort/continue script exists
def script_exists? *args
  File.exist? script_path(*args)
end


# Returns the path to the abort or continue script for the given command
def script_path command, action
  action = action == 'continue' ? '' : "_#{action}"
  "/tmp/#{command.gsub(' ', '-')}#{action}_#{Dir.pwd.gsub '/', '_'}"
end
