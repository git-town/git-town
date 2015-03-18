# Makes finding Git commits easier
class CommitsFinder

  def initialize commit_attributes = [:message]
    @commit_attributes = commit_attributes.map(&:upcase)

    # The currently known commits
    @commits = {}
  end


  # Adds all commits in the current repo
  def add_commits_in_current_repo
    existing_branches.each do |branch_name|
      add_commits_in_branch branch_name
    end
    self
  end


  # Adds the given commit to the list of known commits
  # rubocop:disable MethodLength
  # rubocop:disable AbcSize
  def add_commit sha:, message:, branch_name:
    local_branch_name = local_branch_name branch_name
    @commits[local_branch_name] ||= {}
    if @commits[local_branch_name].key? sha
      # We already have this commit in a different location --> just append the location to the existing commit
      @commits[local_branch_name][sha]['LOCATION'] << branch_location(branch_name)
      return
    end

    commit_data = {
      'BRANCH' => local_branch_name,
      'LOCATION' => [branch_location(branch_name)],
      'MESSAGE' => message
    }
    if @commit_attributes.include? 'FILE NAME'
      filenames = committed_files sha
      commit_data['FILE NAME'] = filenames.join(' and ')
    end
    if @commit_attributes.include? 'FILE CONTENT'
      if filenames.size == 1
        commit_data['FILE CONTENT'] = content_of file: filenames[0], for_sha: sha
      else
        fail 'Cannot verify file content for multiple files'
      end
    end
    @commits[local_branch_name][sha] = commit_data
  end
  # rubocop:enable MethodLength
  # rubocop:enable AbcSize


  # Adds all commits in the given branch
  def add_commits_in_branch branch_name
    array_output_of("git log #{branch_name} --format='%h|%s|%ae' --topo-order --reverse").each do |commit|
      sha, message = commit.split('|')
      next if message == 'Initial commit'
      add_commit sha: sha, message: message, branch_name: branch_name
    end
  end


  # Returns the currently known commits as a Cucumber compatible table
  def to_table
    result = CucumberTableBuilder.new @commit_attributes
    main_commits = @commits.delete 'main'
    main_commits.keys.each do |sha|
      main_commits[sha]['LOCATION'] = main_commits[sha]['LOCATION'].join ' and '
      result.add main_commits[sha].values
    end
    @commits.keys.each do |branch_name|
      @commits[branch_name].keys.each do |sha|
        @commits[branch_name][sha]['LOCATION'] = @commits[branch_name][sha]['LOCATION'].join ' and '
        result.add @commits[branch_name][sha].values
      end
    end
    result.table
  end

end


# Makes it easy to build DRY Cucumber compatible tables
class CucumberTableBuilder

  attr_reader :table

  def initialize headers
    @headers = headers

    @table = [headers]

    # The previously added values
    @previous_values = nil
  end

  def add values
    @table << dry_up(values)
    @previous_values = values
  end

  # Dries up the given values based on what came before in the table
  def dry_up values
    return values unless @previous_values
    result = values.clone
    @previous_values.each_with_index do |previous_value, i|
      if @headers[i] != 'FILE NAME' && values[i] == previous_value
        result[i] = ''
      end
    end
    result
  end

end


# Returns the array of the file names committed for the supplied sha
def committed_files sha
  array_output_of "git diff-tree --no-commit-id --name-only -r #{sha}"
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
  in_secondary_repository do
    create_local_commit commit_data.merge(pull: true, push: true)
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
    commit_data.reverse_merge!(default_commit_attributes)
    create_commit commit_data
  end
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
  output_of "git reflog --grep-reflog='commit: #{commit_message.strip}' --format='%H'"
end


# Verifies the commits in the repository
def verify_commits expected_commits
  actual_commits = CommitsFinder.new(expected_commits.headers).add_commits_in_current_repo.to_table
  expected_commits.diff! actual_commits
end
