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


def create_local_commit options
  run "git checkout #{options[:branch]}"
  File.write options[:file_name], options[:file_content]
  run "git add #{options[:file_name]}"
  run "git commit -m '#{options[:commit_message]}'"
end


# Creates the commits described in the given Cucumber table
#
# The following table columns are supported:
# * branch: branch name
# * location (optional): branch location
#     - local: commit is created in the locally checked out developer repository only
#     - remote: commit is created in the remote repository only
#     - upstream: commit is created in the upstream repository only
#   or any combination of them.
#   Defaults to 'local and remote'
# * message (optional): commit message
# * file name (optional): name of the file to be committed
# * file content (optional): content of the file to be committed
def create_commits commits_table
  current_branch = run('git rev-parse --abbrev-ref HEAD')[:out]

  commits_table.hashes.each do |commit_data|

    # Parse all the given options and augment with default values
    options = {
      file_name: commit_data.delete('file name') { 'default file name' },
      file_content: commit_data.delete('file content') { 'default file content' },
      commit_message: commit_data.delete('message') { 'default commit message' },
      commit_location: commit_data.delete('location'){%i[local remote]},
      branch: commit_data.delete('branch') { current_branch }
    }
    if options[:commit_location].is_a? String
      options[:commit_location] = [options[:commit_location].to_sym]
    end

    # Make sure we understood all commit data
    if commit_data != {}
      raise "Unused commit specifiers: #{commit_data}"
    end

    # Create commits
    if options[:commit_location].delete :local
      create_local_commit options
    end
    if options[:commit_location].delete :remote
      at_path coworker_repository_path do
        run 'git pull'
        create_local_commit options
        run 'git push'
      end
    end
    if options[:commit_location].delete :upstream
      at_path upstream_local_repository_path do
        create_local_commit options
        run 'git push'
      end
    end

    # Make sure we understood all commit data
    if options[:commit_location] != []
      raise "Unused commit location: #{options[:commit_location]}"
    end
  end

  # Go back to the branch that was checked out initially
  run "git checkout #{current_branch}"
end


# Returns whether the given branch name is simple ('feature')
# or not ('remotes/origin/feature')
def is_simple_branch_name branch_name
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

      # Convert simple remote branches ('feature')
      # into their full branch name ('remotes/origin/feature')
      if location == 'remote' && is_simple_branch_name(commit_data_clone[:branch])
        commit_data_clone[:branch] = "remotes/origin/#{commit_data_clone[:branch]}"
      end
      commit_data_clone
    end
  end.flatten

  at_path repository_path do
    expect(commits_in_repo).to match_commits expected_commits
  end
end
