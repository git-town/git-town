def create_local_commit options
  run "git checkout #{options[:branch]}"
  File.write options[:file_name], options[:file_content]
  run "git add #{options[:file_name]}"
  run "git commit -m '#{options[:commit_message]}'"
end


# Returns all commits in all local branches
def actual_commits

  # Save the current branch in order to restore it later
  current_branch = run("git rev-parse --abbrev-ref HEAD")[:out]

  # Get local commits
  actual_commits = existing_local_branches.map do |local_branch_name|
    run "git checkout #{local_branch_name}", allow_failures: true
    commits_in_branch(local_branch_name)
  end.flatten

  # Go back to the branch that was checked out initially
  run "git checkout #{current_branch}", allow_failures: true

  actual_commits
end


# Returns the commits in the currently checked out branch
def commits_in_branch branch_name
  result = run("git log --oneline").fetch(:out)
                                   .split("\n")
                                   .map{|c| { hash: c.slice(0,6),
                                              message: c.slice(8,500),
                                              branch: branch_name }}

  # Remove the root commit
  result.select!{|commit| commit[:message] != 'Initial commit'}

  # Add the affected files to the commits
  result.each do |commit|
    commit[:files] = run("git diff-tree --no-commit-id --name-only -r #{commit[:hash]}").fetch(:out)
                                                                                        .split("\n")
  end

  # Remove the hashes from the commits
  result.each{|c| c.delete(:hash)}

  result
end
