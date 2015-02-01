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
  existing_local_branches.map do |branch|
    files_in(branch: branch).map do |file|
      content = content_of file: file, for_sha: branch
      { branch: branch, name: file, content: content }
    end
  end.flatten
end


def uncommitted_files
  array_output_of "git status --porcelain | awk '{print $2}'"
end


def verify_files files_array
  expected_files = files_array.map do |file_data|
    file_data.symbolize_keys_deep!
    Kappamaki.from_sentence(file_data.delete :files).map do |file|
      content = content_of file: file, for_sha: file_data[:branch]
      file_data.clone.reverse_merge name: file, content: content
    end
  end.flatten

  actual_files = files_in_branches

  expect(actual_files).to match_array expected_files
end
