def format_piece piece, index, column_sizes
  if index == 0
    ' ' * column_sizes[index]
  elsif index == column_sizes.length - 1
    "\n"
  else
    " #{piece.strip} ".ljust column_sizes[index]
  end
end


def format_table lines
  sizes = table_column_sizes lines

  lines.each_with_index.map do |line, line_index|
    line.upcase! if line_index == 0

    line.split('|').each_with_index.map do |piece, index|
      format_piece piece, index, sizes
    end.join('|')
  end
end


def line_column_sizes line
  pieces = line.split('|')
  pieces.each_with_index.map do |piece, index|
    if index == 0 || index == pieces.length - 1
      piece.length
    else
      piece.strip!
      piece.length + 2 # one space padding on each side
    end
  end
end


def table_column_sizes lines
  sizes_by_row = lines.map { |line| line_column_sizes line }
  column_sizes = sizes_by_row[0].zip(*sizes_by_row[1..-1])
  column_sizes.each_with_index.map do |column_size, i|
    i == 0 ? column_size.min : column_size.max
  end
end


perform_clean = ARGV.length == 0


Dir.glob('./features/**/*.feature').each do |filename|
  out_lines = []
  table_lines = []

  File.open(filename) do |f|
    f.each_line do |line|
      if line.include? '|'
        table_lines << line
      else
        unless table_lines.empty?
          out_lines += format_table(table_lines)
          table_lines = []
        end

        out_lines << line
      end
    end
  end

  out_lines += format_table(table_lines) unless table_lines.empty?

  if perform_clean
    File.open(filename, 'w') { |f| f.write(out_lines.join('')) }
  elsif File.read(filename) != out_lines.join('')
    puts ''
    puts "#{filename} has an ill-formatted cucumber table"
    puts 'Run "bin/clean_cucumber" to fix.'
    puts ''
    exit 1
  end
end
