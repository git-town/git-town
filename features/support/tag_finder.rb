# Finds Git tags, returns them as TagList instances
class TagFinder

  # Returns all tags (local + remote)
  def self.all_tags
    TagList.new
      .add_many(location: 'local', tags: tags_in(:developer))
      .add_many(location: 'remote', tags: tags_in(:origin))
      .to_table
  end


  # Returns the names of the Git tags in the given repository
  def self.tags_in repo
    in_repository repo do
      array_output_of('git tag')
    end
  end
  private_class_method :tags_in

end
