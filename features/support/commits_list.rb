# Encapsulates building up a list of commits
#
# Commit lists contain very specific data about commits.
# They are also sorted in distinct ways.
class CommitsList

  # commit_attributes - the attributes we want to gather for the commits found
  #                     The default value is the minimalistic set
  #                     for checking if commits exist at all
  #
  # Examples
  #
  #   CommitsList.new [:branch, :location, :message, :file_name, :file_content, :author]
  def initialize commit_attributes = [:message]
    @commit_attributes = commit_attributes.map(&:upcase)

    # The currently known commits
    @commits = {}
  end


  # Adds the given commit to this CommitsList instance
  #
  # rubocop:disable MethodLength
  # rubocop:disable AbcSize
  def add sha:, message:, branch_name:, author:
    local_branch_name = local_branch_name branch_name
    @commits[local_branch_name] ||= {}
    if @commits[local_branch_name].key? sha
      # We already have this commit in a different location --> just append the location to the existing commit
      @commits[local_branch_name][sha]['LOCATION'] << branch_location(branch_name)
    else
      # We don't have this commit in another location --> create an entry for this commit
      @commits[local_branch_name][sha] = commit_data local_branch_name: local_branch_name,
                                                     branch_name: branch_name,
                                                     message: message,
                                                     sha: sha,
                                                     author: author
    end
  end
  # rubocop:enable MethodLength
  # rubocop:enable AbcSize


  # Returns whether this CommitsList instance has found any commits so far
  def empty?
    @commits.empty?
  end


  # Returns the currently known commits as a Cucumber compatible array
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
  # rubocop:enable AbcSize
  # rubocop:enable MethodLength


  private

  # Returns whether this CommitsList instance is looking for
  # the given commit attribute
  def attribute? attribute_name
    @commit_attributes.include? attribute_name
  end


  # Returns the internally used data structure for the commit with the given data
  #
  # rubocop:disable MethodLength
  def commit_data local_branch_name:, branch_name:, message:, sha:, author:
    result = {
      'BRANCH' => local_branch_name,
      'LOCATION' => [branch_location(branch_name)],
      'MESSAGE' => message
    }
    if attribute? 'FILE NAME'
      filenames = committed_files sha
      result['FILE NAME'] = filenames[0] || ''
    end
    if attribute? 'FILE CONTENT'
      if filenames.size == 1
        result['FILE CONTENT'] = content_of file: filenames[0], for_sha: sha
      else
        fail 'Cannot verify file content for multiple files'
      end
    end
    result['AUTHOR'] = author if attribute? 'AUTHOR'
    result
  end
  # rubocop:enable MethodLength

end
