# Returns the branch hierarchy information that is currently configured in the current repository
def configured_branch_hierarchy_information ignore_errors: false
  result = Mortadella.new headers: %w(BRANCH PARENT)
  array_output_of("git config --get-regexp '^git-town\\.branches\\.parent\\.'",
                  ignore_errors: ignore_errors).sort.each do |row|
    result << row.scan(/^git-town\.branches\.parent\.(.+) (.+)$/)[0]
  end
  result
end


# Stores the currently configured branch hierarchy metadata for later
def store_branch_hierarchy_metadata
  @branch_hierarchy_metadata = configured_branch_hierarchy_information
end
