# Makes it easy to build DRY Cucumber-compatible tables
class CucumberTableBuilder

  attr_reader :table


  def initialize headers:, dry: []
    @headers = headers

    @dry = dry

    # The resulting Cucumber-compatible table structure
    @table = [headers]

    # The previously added row
    @previous_row = nil
  end


  # Adds the given row to the table
  def << row
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
      if @dry.include?(@headers[i]) && row[i] == previous_value
        result[i] = ''
      else
        break
      end
    end
    result
  end

end
