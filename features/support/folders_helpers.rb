# Returns the full path to the folder with the given name inside the current
# Git workspace
def git_folder folder_name
  "#{git_root_folder}/#{folder_name}"
end


def git_root_folder
  output_of 'git rev-parse --show-toplevel'
end
