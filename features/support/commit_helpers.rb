# frozen_string_literal: true
# Returns the array of the file names committed for the supplied sha
def committed_files sha
  array_output_of "git diff-tree --no-commit-id --name-only -r #{sha}"
end


# Returns all commits in the current repository as a data table
def commits_in_current_repo commit_fields
  CommitListBuilder.new(commit_fields).add_commits_in_current_repo.to_table
end


# Creates a new commit with the given properties.
#
# Parameter is a Cucumber table line
def create_local_commit author: nil, branch:, file_name:, file_content:, message:, push: false
  on_branch(branch) do
    if (folder_name = File.dirname file_name) != '.'
      Dir.mkdir folder_name
    end
    File.write file_name, file_content
    run "git add '#{file_name}'"
    run "git commit -m '#{message}' #{"--author='#{author}'" if author}"
    run 'git push' if push
  end
end


def create_remote_commit commit_data
  in_secondary_repository do
    run 'git fetch --prune'
    create_local_commit commit_data.merge(push: true)
  end
end


def create_upstream_commit commit_data
  in_repository :upstream_developer do
    create_local_commit commit_data.merge(push: true)
  end
end


def create_commit commit_data
  location = Kappamaki.from_sentence commit_data.delete(:location)

  case location
  when %w(local) then create_local_commit commit_data
  when %w(remote) then create_remote_commit commit_data
  when %w(local remote) then create_local_commit commit_data.merge(push: true)
  when %w(upstream) then create_upstream_commit commit_data
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
  normalize_commit_data commits_array
  commits_array.each do |commit_data|
    create_commit add_default_commit_data(commit_data)
  end
end


# Returns the SHA for the commit with the given message
def commit_sha commit_message
  output_of "git log --grep='#{commit_message}' --format='%h' -1"
end


# Returns a commit_data structure consisting of
# the given commit_data structure with default values added
def add_default_commit_data commit_data
  file_name = "default_file_name_#{SecureRandom.urlsafe_base64}"
  commit_data.clone.reverse_merge(
    file_name: file_name,
    message: 'default commit message',
    location: 'local and remote',
    branch: current_branch_name,
    file_content: default_file_content_for(commit_data[:file_name] || file_name)
  )
end


# Returns the file content that is used if no file content is provided by the user
def default_file_content_for file_name
  "#{file_name} content"
end


def group_expected_commits_by_branch commits_array
  out = {}

  commits_array.each do |commit_data|
    branch = commit_data.delete(:branch)
    out[branch] ||= []
    out[branch] << commit_data
  end

  out
end


# Normalize commits_array by converting all keys to symbols and
# filling in any data implied from the previous commit
def normalize_commit_data commits_array
  commits_array = commits_array.each(&:symbolize_keys_deep!)
  commits_array.each_cons(2) do |previous_commit_data, commit_data|
    commit_data.default_blank! previous_commit_data.subhash(:branch, :location)
  end
end


# Normalize commit_data by parsing the files and location
# Returns an array of commit_data
def normalize_expected_commit_data commit_data
  # Convert file string list into real array
  commit_data[:file_name] = Kappamaki.from_sentence commit_data[:file_name]

  # Create individual expected commits for each location provided
  Kappamaki.from_sentence(commit_data.delete(:location)).map do |location|
    branch = branch_name_for_location location, commit_data[:branch]
    commit_data.clone.merge branch: branch
  end
end


# Returns an array of length count with the shas of the most recent commits
def recent_commit_shas count
  array_output_of("git rev-list HEAD -n #{count}").reverse
end


# Returns the SHA of the commit with the given message
def sha commit_message
  cmd = "git reflog --grep-reflog='commit"
  cmd += ' (initial)' if commit_message == 'Initial commit'
  cmd += ": #{commit_message.strip}' --format='%H'"
  result = ''
  %w(developer developer_secondary).each do |user|
    next unless result.empty?
    in_repository(user) { result = output_of cmd }
  end
  result
end


# Verifies the commits in the repository
def verify_commits expected_commits
  expected_commits.diff! commits_in_current_repo(expected_commits.headers)
end
