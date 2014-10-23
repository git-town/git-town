# Returns the commits in the current directory
def commits_in_repo
  current_branch = run('git rev-parse --abbrev-ref HEAD')[:out]

  result = []
  existing_local_branches.map do |local_branch_name|
    run "git checkout #{local_branch_name}"
    commits = local_commits.each do |commit|
      commit[:location] = 'local'
      commit[:branch] = local_branch_name
    end
    result.concat commits
  end

  existing_remote_branches.map do |remote_branch_name|
    run "git checkout #{remote_branch_name}"
    commits = local_commits.each do |commit|
      commit[:location] = 'remote'
      commit[:branch] = remote_branch_name
    end
    result.concat commits
  end

  run "git checkout #{current_branch}"
  result
end


# Creates a new commit with the given properties.
#
# Parameter is a Cucumber table line
def create_local_commit branch:, file_name:, file_content:, message:
  run "git checkout #{branch}"
  File.write file_name, file_content
  run "git add '#{file_name}'"
  run "git commit -m '#{message}'"
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
  current_branch = run('git rev-parse --abbrev-ref HEAD')[:out]

  commits_array = [commits_array] if commits_array.is_a? Hash
  commits_array.each do |commit_data|
    symbolize_keys_deep! commit_data

    # Augment the commit data with default values
    commit_data.reverse_merge!({ file_name: "default file name #{SecureRandom.urlsafe_base64}",
                                 file_content: 'default file content',
                                 message: 'default commit message',
                                 location: 'local and remote',
                                 branch: current_branch })

    # Create the commits
    case (location = Kappamaki.from_sentence commit_data.delete(:location))
    when %w[local]
      create_local_commit commit_data
    when %w[remote]
      at_path coworker_repository_path do
        run 'git pull'
        create_local_commit commit_data
        run 'git push'
      end
    when %w[local remote]
      create_local_commit commit_data
      run 'git push'
    when %w[upstream]
      at_path upstream_local_repository_path do
        create_local_commit commit_data
        run 'git push'
      end
    else
      raise "Unknown commit location: #{location}"
    end
  end

  # Go back to the branch that was checked out initially
  run "git checkout #{current_branch}"
end


# Returns whether the given branch name is simple ('feature')
# or not ('remotes/origin/feature')
def is_local_branch_name? branch_name
  /^[^\/]+$/.match branch_name
end

# Returns the commits in the currently checked out branch
def local_commits
  result = run("git log --oneline").fetch(:out)
                                   .split("\n")
                                   .map{|c| { hash: c.slice(0,6),
                                              message: c.slice(8,500) }}

  # Remove the root commit
  result.select!{|commit| commit[:message] != 'Initial commit'}

  # Add the affected files to the commits
  result.each do |commit|
    commit[:files] = run("git diff-tree --no-commit-id --name-only -r #{commit[:hash]}").fetch(:out)
                                                                                        .split("\n")
  end

  # Remove the hashes from the commits
  result.each{|commit| commit.delete(:hash)}

  result
end


# Verifies that the commits in the repository at the given path
# are similar to the expected commits in the given Cucumber table
def verify_commits commits_table:, repository_path:
  expected_commits = commits_table.hashes.map do |commit_data|
    symbolize_keys_deep! commit_data

    # Convert file string list into real array
    commit_data[:files] = Kappamaki.from_sentence commit_data[:files]

    # Create individual expected commits for each location provided
    Kappamaki.from_sentence(commit_data[:location]).map do |location|
      commit_data_clone = commit_data.clone
      commit_data_clone[:location] = location

      # Convert a local branch name ('feature')
      # into its corresponding tracking branch name ('remotes/origin/feature')
      if location == 'remote' && is_local_branch_name?(commit_data_clone[:branch])
        commit_data_clone[:branch] = "remotes/origin/#{commit_data_clone[:branch]}"
      end
      commit_data_clone
    end
  end.flatten

  at_path repository_path do
    expect(commits_in_repo).to match_commits expected_commits
  end
end
