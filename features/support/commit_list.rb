# Encapsulates building up a list of commits
#
# Commit lists contain very specific data about commits.
# They are also sorted in distinct ways.
class CommitListBuilder

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
  def add sha:, message:, branch_name:, author:
    local_branch_name = local_branch_name branch_name
    @commits[local_branch_name] ||= {}
    if @commits[local_branch_name].key? sha
      # We already have this commit in a different location --> just append the location to the existing commit
      @commits[local_branch_name][sha]['LOCATION'] << branch_location(branch_name)
      return
    end

    @commits[local_branch_name][sha] = commit_data local_branch_name: local_branch_name, branch_name: branch_name, message: message, sha: sha, author: author
  end
  # rubocop:enable MethodLength
  # rubocop:enable AbcSize


  # Returns whether this CommitsFinder instance has found any commits so far
  def empty?
    @commits.empty?
  end


  # Returns the currently found commits as a Cucumber compatible table
  #
  # rubocop:disable MethodLength
  # rubocop:disable AbcSize
  def to_table
    result = CucumberTableBuilder.new headers: @commit_attributes, dry: %w(BRANCH LOCATION)
    main_commits = @commits.delete 'main'
    main_commits.try(:keys).try(:each) do |sha|
      main_commits[sha]['LOCATION'] = main_commits[sha]['LOCATION'].to_sentence
      result << main_commits[sha].values
    end
    @commits.values.each do |branch_commits|
      branch_commits.values.each do |commit|
        commit['LOCATION'] = commit['LOCATION'].to_sentence
        result << commit.values
      end
    end
    result.table
  end
  # rubocop:enable MethodLength
  # rubocop:enable AbcSize


private

  # Returns whether this CommitsFinder instance is looking for
  # the given commit attribute
  def attribute? attribute_name
    @commit_attributes.include? attribute_name
  end


  def commit_data local_branch_name:, branch_name:, message:, sha:, author:
    commit_data = {
      'BRANCH' => local_branch_name,
      'LOCATION' => [branch_location(branch_name)],
      'MESSAGE' => message
    }
    if attribute? 'FILE NAME'
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
    commit_data
  end
end
