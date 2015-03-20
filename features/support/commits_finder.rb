# Makes finding Git commits easier
class CommitsFinder

  # commit_attributes - the attributes we want to gather for the commits found
  #                     The default value is the minimalistic set
  #                     for checking if commits exist at all
  #
  # Examples
  #
  #   CommitsFinder.new [:branch, :location, :message, :file_name, :file_content, :author]
  def initialize commit_attributes = [:message]
    @commit_attributes = commit_attributes.map(&:upcase)

    # The currently known commits
    @commits = {}
  end


  # Adds the given commit to this CommitsFinder instance
  #
  # rubocop:disable MethodLength
  # rubocop:disable AbcSize
  def add_commit sha:, message:, branch_name:, author:
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
      commit_data['FILE NAME'] = filenames.to_sentence
    end
    if @commit_attributes.include? 'FILE CONTENT'
      if filenames.size == 1
        commit_data['FILE CONTENT'] = content_of file: filenames[0], for_sha: sha
      else
        fail 'Cannot verify file content for multiple files'
      end
    end
    commit_data['AUTHOR'] = author if @commit_attributes.include? 'AUTHOR'
    @commits[local_branch_name][sha] = commit_data
  end
  # rubocop:enable MethodLength
  # rubocop:enable AbcSize


  # Adds all commits in the given branch to this CommitsFinder instance
  def add_commits_in_branch branch_name
    array_output_of("git log #{branch_name} --format='%h|%s|%ae' --topo-order --reverse").each do |commit|
      sha, message, author = commit.split('|')
      next if message == 'Initial commit'
      add_commit sha: sha, message: message, branch_name: branch_name, author: author
    end
  end


  # Adds all commits in the current repo to this CommitsFinder instance
  def add_commits_in_current_repo
    existing_branches.each do |branch_name|
      add_commits_in_branch branch_name
    end
    self
  end


  # Returns whether this CommitsFinder instance has found any commits so far
  def empty?
    @commits.empty?
  end


  # Returns the currently found commits as a Cucumber compatible table
  #
  # rubocop:disable MethodLength
  # rubocop:disable AbcSize
  def to_table
    result = CucumberTableBuilder.new @commit_attributes
    main_commits = @commits.delete 'main'
    main_commits.try(:keys).try(:each) do |sha|
      main_commits[sha]['LOCATION'] = main_commits[sha]['LOCATION'].to_sentence
      result.add_row main_commits[sha].values
    end
    @commits.keys.each do |branch_name|
      @commits[branch_name].keys.each do |sha|
        @commits[branch_name][sha]['LOCATION'] = @commits[branch_name][sha]['LOCATION'].to_sentence
        result.add_row @commits[branch_name][sha].values
      end
    end
    result.table
  end
  # rubocop:enable MethodLength
  # rubocop:enable AbcSize

end
