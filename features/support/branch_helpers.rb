# Returns the branch name for the given location
def branch_name_for_location location, branch
  case location
  when 'local', 'coworker' then branch
  when 'remote' then "origin/#{branch}"
  when 'upstream' then "upstream/#{branch}"
  else fail "Unknown location: #{location}"
  end
end


# Returns the location of the branch with the given name
#
# 'foo'          --> 'local'
# 'origin/foo'   --> 'remote'
# 'upstream/foo' --> 'upstream'
def branch_location branch_name
  case
  when branch_name.start_with?('origin/') then 'remote'
  when branch_name.start_with?('upstream/') then 'upstream'
  else 'local'
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
def create_branch branch_name, remote: true, start_point: 'main'
  run "git checkout -b #{branch_name} #{start_point}"
  run "git push -u origin #{branch_name}" if remote
end


# Returns the name of the branch that is currently checked out
def current_branch_name
  output_of 'git rev-parse --abbrev-ref HEAD'
end


# Returns the names of the existing feature branches
def existing_branches order: :alphabetically
  existing_local_branches(order: order) + existing_remote_branches
end


# Returns the names of all existing local branches.
def existing_local_branches order: :alphabetically
  result = array_output_of "git branch | tr -d '*'"
  if order == :main_first
    main_branch = result.delete 'main'
    result = [main_branch].concat result if main_branch
  end
  result
end


# Returns the names of all existing remote branches.
def existing_remote_branches
  remote_branches = array_output_of 'git branch -r'
  remote_branches.reject { |b| b.include?('HEAD') }
end


# Returns the name of the given branch if it was local
#
# 'foo'          --> 'foo'
# 'origin/foo'   --> 'foo'
# 'upstream/foo' --> 'foo'
def local_branch_name branch_name
  branch_name.sub(/.+\//, '')
end


# Returns the given branch name in a format that is compatible with Git config
def normalize_branch_name branch_name
  branch_name.gsub '_', '-'
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


def set_parent_branch branch:, parent:, parents:
  run "git config git-town.branches.parent.#{normalize_branch_name branch} #{parent}"
  run "git config git-town.branches.parents.#{normalize_branch_name branch} #{parents}"
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
