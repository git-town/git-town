# Returns the branch name for the given location
def branch_name_for_location location, branch
  case location
  when 'local', 'coworker' then branch
  when 'remote' then "origin/#{branch}"
  when 'upstream' then "upstream/#{branch}"
  else fail "Unknown location: #{location}"
  end
end


# Returns the branches for the given repository
def branches_for_repository repository
  case repository
  when 'local' then existing_local_branches
  when 'remote' then existing_remote_branches
  when 'coworker' then in_repository(:coworker) { existing_local_branches }
  else fail "Unknown repository: #{repository}"
  end
end


# Creates and pushes a branch
def create_branch branch_name, remote: true
  run "git checkout -b #{branch_name} main"
  run "git push -u origin #{branch_name}" if remote
end


# Returns the name of the branch that is currently checked out
def current_branch_name
  output_of 'git rev-parse --abbrev-ref HEAD'
end


# Returns the names of the existing feature branches
def existing_branches
  existing_local_branches + existing_remote_branches
end


# Returns the names of all existing local branches.
def existing_local_branches
  array_output_of "git branch | tr -d '*'"
end


# Returns the names of all existing remote branches.
def existing_remote_branches
  remote_branches = array_output_of 'git branch -r'
  remote_branches.reject { |b| b.include?('HEAD') }
end


def number_of_branches_out_of_sync
  integer_output_of 'git branch -vv | grep -o "\[.*\]" | tr -d "[]" | awk "{ print \$2 }" | grep . | wc -l'
end


# Executes the given block after checking out the supplied branch
# then returns to the currently branch
def on_branch branch_name
  original_branch = current_branch_name
  run "git checkout #{branch_name}"
  result = yield
  run "git checkout #{original_branch}"
  result
end


# Returns the SHA of the given branch
def sha_of_branch branch_name
  output_of "git rev-parse #{branch_name}"
end


# Verifies the branches in each repository
def verify_branches branch_table
  branch_table.map_headers!(&:downcase)
  branch_table.hashes.each do |branch_data|
    repository = branch_data['repository']
    expected_branches = Kappamaki.from_sentence branch_data['branches']
    expected_branches.map! { |branch_name| branch_name_for_location repository, branch_name }
    actual_branches = branches_for_repository repository
    expect(actual_branches).to match_array(expected_branches)
  end
end
