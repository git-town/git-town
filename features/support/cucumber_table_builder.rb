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

    # The previously added row
    @previous_row = nil
  end


  # Adds the given row to the table
  def add_row row
    @table << dry_up(row)
    @previous_row = row
  end


  # Returns a dried up version of the given row
  # based on the row that came before in the table
  #
  # In a dried up row, any values that match the previous row are removed,
  # stopping on the first difference
  def dry_up row
    return row unless @previous_row
    result = row.clone
    @previous_row.each_with_index do |previous_value, i|
      if @headers[i] != 'MESSAGE' && @headers[i] != 'FILE NAME' && row[i] == previous_value
        result[i] = ''
      else
        break
      end
    end
    result
  end

end
