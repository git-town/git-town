# Returns the commits in the current directory
def commits_in_repo
  existing_branches.map do |branch_name|
    commits_for_branch branch_name
  end.flatten
end


# Creates a new commit with the given properties.
#
# Parameter is a Cucumber table line
def create_local_commit branch:, file_name:, file_content:, message:, push: false, pull: false
  run 'git fetch --prune' if pull
  on_branch(branch) do
    File.write file_name, file_content
    run "git add '#{file_name}'"
    run "git commit -m '#{message}'"
    run 'git push' if push
  end
end


def create_remote_commit commit_data
  at_path coworker_repository_path do
    create_local_commit commit_data.merge(pull: true, push: true)
  end
end


def create_upstream_commit commit_data
  at_path upstream_local_repository_path do
    create_local_commit commit_data.merge(push: true)
  end
end


def create_commit commit_data
  location = Kappamaki.from_sentence commit_data.delete(:location)

  case location
  when %w(local) then create_local_commit commit_data
  when %w(remote) then create_remote_commit commit_data
  when %w(local remote) then create_local_commit commit_data.merge(push: true)
  when %w(upstream) then  create_upstream_commit commit_data
  else fail "Unknown commit location: #{location}"
  end
end


# Creates the commits described in an array of hashs
#
# The following keys are supported. All of them are optional:
# | column name  | default                | description                                                |
# | branch       | current branch         | name of the branch in which to create the commit           |
# | location     | local and remote       | where to create the commit                                 |
# |              |                        | - local: the locally checked out developer repository only |
# |              |                        | - remote: in the remote repository only                    |
# |              |                        | - upstream: in the upstream repository only                |
# | message      | default commit message | commit message                                             |
# | file name    | default file name      | name of the file to be committed                           |
# | file content | default file content   | content of the file to be committed                        |
def create_commits commits_array
  commits_array = [commits_array] if commits_array.is_a? Hash
  commits_array.each do |commit_data|
    commit_data.keys.each do |key|
      commit_data[key.downcase] = commit_data.delete key
    end
    symbolize_keys_deep! commit_data
    commit_data.reverse_merge!(default_commit_attributes)
    create_commit commit_data
  end
end


# Returns the array of the file names committed for the supplied sha
def committed_files sha
  array_output_of "git diff-tree --no-commit-id --name-only -r #{sha}"
end


# Returns the commits in the currently checked out branch
def commits_for_branch branch_name
  array_output_of("git log #{branch_name} --oneline").map do |commit|
    sha, message = commit.split(' ', 2)

    unless message == 'Initial commit'
      { branch: branch_name, message: message, files: committed_files(sha) }
    end
  end.compact
end


def default_commit_attributes
  {
    file_name: "default file name #{SecureRandom.urlsafe_base64}",
    file_content: 'default file content',
    message: 'default commit message',
    location: 'local and remote',
    branch: current_branch_name
  }
end


# Normalize commit_data by parsing the files and location
# Returns an array of commit_data
def normalize_expected_commit_data commit_data
  commit_data.keys.each do |key|
    commit_data[key.downcase] = commit_data.delete key
  end
  symbolize_keys_deep! commit_data

  # Convert file string list into real array
  commit_data[:files] = Kappamaki.from_sentence commit_data[:files]

  # Create individual expected commits for each location provided
  Kappamaki.from_sentence(commit_data.delete(:location)).map do |location|
    branch = branch_name_for_location location, commit_data[:branch]
    commit_data.clone.merge branch: branch
  end
end


# Returns an array of length count with the shas of the most recent commits
def recent_commit_shas count
  array_output_of("git rev-list HEAD -n #{count}")
end


# Verifies that the commits in the repository at the given path
# are similar to the expected commits in the given Cucumber table
def verify_commits commits_table:, repository_path:
  expected_commits = commits_table.hashes.map do |commit_data|
    normalize_expected_commit_data commit_data
  end.flatten

  at_path repository_path do
    actual_commits = commits_in_repo
    expect(actual_commits).to match_array(expected_commits), -> { commits_diff(actual_commits, expected_commits) }
  end
end
