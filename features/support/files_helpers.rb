# Returns the content of the file with the given name on the given branch
def content_of file:, in_branch:
  output_of "git show #{in_branch}:#{file}"
end


# Provides the files that exist in the given branch
def files_in_branch branch_name
  array_output_of "git ls-tree -r --name-only #{branch_name}"
end


# Returns a table of all files in all branches.
# This is for comparing against expected files in a Cucumber table.
def all_files_in_all_branches except: []
  [].tap do |result|
    existing_local_branches.each do |branch|
      files_in_branch(branch).each do |file_path|
        unless except.include? file_path
          result << { branch: branch, name: file_path, content: content_of(file: file_path, in_branch: branch) }
        end
      end
    end
  end
end


def uncommitted_files
  array_output_of "git status --porcelain | awk '{print $2}'"
end
