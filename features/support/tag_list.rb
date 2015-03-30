# Represents a list of Git tags and their locations
class TagList

  def initialize
    @tags = {}
  end


  # Adds a tag with the given name and location to this list
  def add name:, location:
    @tags[name] ||= []
    @tags[name] << location
  end


  # Returns this tag list as a Mortadella instance
  def to_table
    result = Mortadella.new headers: %w(NAME LOCATION)
    @tags.keys.sort.each do |tag_name|
      result << [tag_name, @tags[tag_name].to_sentence]
    end
    result
  end

end
