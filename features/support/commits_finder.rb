# Makes finding Git commits easier
#
# After initializing with the desired attributes for the commit list,
# run as many #find_* methods as you want to build up your commit list.
# When done, call #to_a to get the results as a Cucumber compatible array.
class CommitsFinder

  # commit_attributes - the attributes we want to gather for the commits found
  #                     The default value is the minimalistic set
  #                     for checking if commits exist at all
  #
  # Examples
  #
  #   CommitsFinder.new [:branch, :location, :message, :file_name, :file_content, :author]
  def initialize commit_attributes = [:message]
    @commits = CommitListBuilder.new commit_attributes
  end


  # Adds all commits in the given branch to this CommitsFinder instance
  def find_commits_in_branch branch_name
    array_output_of("git log #{branch_name} --format='%h|%s|%ae' --topo-order --reverse").each do |commit|
      sha, message, author = commit.split('|')
      next if message == 'Initial commit'
      @commits.add sha: sha, message: message, branch_name: branch_name, author: author
    end
  end


  # Adds all commits in the current repo to this CommitsFinder instance
  def find_commits_in_current_repo
    existing_branches.each do |branch_name|
      find_commits_in_branch branch_name
    end
    self
  end


  # Returns whether this CommitsFinder instance has found any commits so far
  def empty?
    @commits.empty?
  end


  def to_a
    @commits.to_table
  end

end
