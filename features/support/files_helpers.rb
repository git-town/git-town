# Returns the content of the file with the given name on the given branch
def content_of file:, in_branch:
  output_of "git show #{in_branch}:#{file}"
end


# Provides the files that exist in the given branch
def files_in branch:
  array_output_of "git ls-tree -r --name-only #{branch}"
end


# Returns a table of all files in all branches.
# This is for comparing against expected files in a Cucumber table.
def all_files_in_all_branches except: []
  [].tap do |result|
    for branch in existing_local_branches
      for file in files_in branch: branch
        unless except.include? file
          result << { branch: branch,
                      name: file,
                      content: content_of(file: file, in_branch: branch) }
        end
      end
    end
  end
end


def uncommitted_files
  array_output_of "git status --porcelain | awk '{print $2}'"
end
