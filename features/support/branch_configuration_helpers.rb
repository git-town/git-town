# Returns the branch hierarchy information that is currently configured in the current repository
def configured_branch_hierarchy_information
  result = Mortadella.new headers: %w(BRANCH PARENT)
  array_output_of("git config --get-regexp '^git-town\\.branches\\.parent\\.'").sort.each do |row|
    result << row.scan(/^git-town\.branches\.parent\.(.+) (.+)$/)[0]
  end
  result
end
