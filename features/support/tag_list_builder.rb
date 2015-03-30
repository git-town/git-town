# Builds up lists of tags
class TagListBuilder

  def initialize
    @tag_list = TagList.new
  end


  # Adds all the tags in the given repo location
  #
  # Example:
  #   TagListBuilder.new.add_tags 'local'
  #   TagListBuilder.new.add_tags ['local', 'remote']
  def add_tags location
    send("#{location}_tags").each do |tag_name|
      @tag_list.add name: tag_name, location: location
    end
    self
  end


  # Returns a Mortadella instance containing the currently built tag list
  def to_table
    @tag_list.to_table
  end



  private

  # Returns all local tags
  def local_tags
    array_output_of 'git tag'
  end


  # Returns all tags that exist in the remote repo
  def remote_tags
    in_repository :origin do
      local_tags
    end
  end

end
