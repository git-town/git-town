# Returns the content of the file with the given name on the given branch
def content_of file:, for_sha:
  output_of "git show #{for_sha}:#{file}"
end


# Provides the files that exist in the given branch
def files_in branch:
  array_output_of("git ls-tree -r --name-only #{branch}")
end


# Returns a table of all files in all branches.
# This is for comparing against expected files in a Cucumber table.
def files_in_branches
  result = [%w(BRANCH NAME CONTENT)]
  existing_local_branches(order: :main_first).each do |branch|
    files_in(branch: branch).each do |file|
      result << [branch, file, content_of(file: file, for_sha: branch)]
    end
  end
  result
end


def create_uncommitted_file options = {}
  options.reverse_merge!(default_uncommitted_file_attributes)
  IO.write options[:file_name], options[:file_content]
  options[:file_name]
end


def default_uncommitted_file_attributes
  {
    file_name: "uncommitted_file",
    file_content: 'default uncommitted content'
  }
end


def verify_uncommitted_file options = {}
  if options.key? :file_name
    options.reverse_merge!(default_uncommitted_file_attributes)
    expect(uncommitted_files).to include options[:file_name]
    expect(IO.read options[:file_name]).to eql options[:file_content]
  else
    fail 'Cannot verify uncommitted file without file_name'
  end
end


def uncommitted_files
  array_output_of "git status --porcelain | awk '{print $2}'"
end
