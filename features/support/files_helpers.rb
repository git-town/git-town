# Returns the content of the file with the given name on the given branch
def content_of file:, in_branch:
  output_of "git show #{in_branch}:#{file}"
end


# Provides the files that exist in the given branch
def files_in branch:
  array_output_of("git ls-tree -r --name-only #{branch}") - ['.gitignore']
end


# Returns a table of all files in all branches.
# This is for comparing against expected files in a Cucumber table.
def files_in_branches
  existing_local_branches.map do |branch|
    files_in(branch: branch).map do |file|
      content = content_of file: file, in_branch: branch
      { branch: branch, name: file, content: content }
    end
  end.flatten
end


def uncommitted_files
  array_output_of "git status --porcelain | awk '{print $2}'"
end
