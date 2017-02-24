# frozen_string_literal: true
# Represents a sorted list of Git tags and their locations
class TagList

  def initialize
    @tags = {}
  end


  # Adds the Git tag with the given name and location to this tag list
  def add name:, location:
    @tags[name] ||= []
    @tags[name] << location
  end


  # Adds the given list of Git tags from the given location to this tag list
  #
  # This call is chainable
  def add_many location:, tags:
    tags.each do |tag|
      add name: tag, location: location
    end
    self
  end


  # Returns this tag list as a Cucumber-compatible table
  def to_table
    mortadella = Mortadella::Horizontal.new headers: %w(NAME LOCATION)
    @tags.keys.sort.each do |tag_name|
      mortadella << [tag_name, @tags[tag_name].to_sentence]
    end
    mortadella.table
  end

end
