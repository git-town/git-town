# Returns the content of the file with the given name on the given branch
def content_of file:, for_sha:
  result = output_of "git show #{for_sha}:#{file}"
  result = '' if result == default_file_content_for(file)
  result
end


# Provides the files that exist in the given branch
def files_in branch:
  array_output_of("git ls-tree -r --name-only #{branch}")
end


# Returns a table of all files in all branches.
# This is for comparing against expected files in a Cucumber table.
def files_in_branches
  result = Mortadella::Horizontal.new headers: %w(BRANCH NAME CONTENT)
  existing_local_branches(order: :main_first).each do |branch|
    files_in(branch: branch).each do |file|
      result << [branch, file, content_of(file: file, for_sha: branch)]
    end
  end
  result.table
end


def create_uncommitted_file options = {}
  options.reverse_merge! DEFAULT_UNCOMMITTED_FILE_ATTRIBUTES
  dirname = File.dirname options[:name]
  FileUtils.mkdir_p(dirname) unless File.directory? dirname
  IO.write options[:name], options[:content]
  options[:name]
end


DEFAULT_UNCOMMITTED_FILE_ATTRIBUTES = {
  name: 'uncommitted_file',
  content: 'default uncommitted content'
}.freeze


def verify_uncommitted_file options
  options.reverse_merge! DEFAULT_UNCOMMITTED_FILE_ATTRIBUTES
  expect(uncommitted_files).to include options[:name]
  expect(IO.read options[:name]).to eql options[:content]
end


def uncommitted_files
  array_output_of "git status --porcelain --untracked-files=all | awk '{print $2}'"
end
