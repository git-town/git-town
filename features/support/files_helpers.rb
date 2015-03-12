# Returns the content of the file with the given name on the given branch
def content_of file:, for_sha:
  output_of "git show #{for_sha}:#{file}"
end


# Provides the files that exist in the given branch
def files_in branch:
  array_output_of("git ls-tree -r --name-only #{branch}")
end


def files_in_branches
  result = [%w(BRANCH NAME CONTENT)]
  existing_local_branches.each do |branch|
    files_in(branch: branch).map do |file|
      content = content_of file: file, for_sha: branch
      result << [branch, file, content]
    end
  end
  result
end


def uncommitted_files
  array_output_of "git status --porcelain | awk '{print $2}'"
end
