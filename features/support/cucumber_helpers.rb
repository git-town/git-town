# Makes it easy to build DRY Cucumber-compatible tables
#
# Example
#
#   table_builder = CucumberTableBuilder.new ['AUTHOR', 'MESSAGE', 'FILES']
#   table_builder.add ['Jay', 'hello', 'jay.txt']
#   table_builder.add ['Jane', 'hi', 'jane.txt']
#   table_builder.to_table
#
#   =>  [ ['AUTHOR', 'MESSAGE', 'FILES'],
#         ['Jay', 'hello', 'jay.txt'],
#         ['Jane', 'hi', 'jane.txt'] ]
class CucumberTableBuilder

  attr_reader :table


  def initialize headers
    @headers = headers

    # The resulting Cucumber-compatible table structure
    @table = [headers]

    # The previously added row of values
    @previous_values = nil
  end


  # Adds the given row of values to the table
  def add_row values
    @table << dry_up(values)
    @previous_values = values
  end


  # Dries up the given values based on what came before in the table
  #
  # rubocop:disable MethodLength
  def dry_up values
    return values unless @previous_values
    result = values.clone
    previous_column_empty = true   # indicates whether the data at the previous value of i was
    @previous_values.each_with_index do |previous_value, i|
      if @headers[i] != 'MESSAGE' && @headers[i] != 'FILE NAME' && values[i] == previous_value && previous_column_empty
        result[i] = ''
        previous_column_empty = true
      else
        previous_column_empty = false
      end
    end
    result
  end
  # rubocop:enable MethodLength

end
