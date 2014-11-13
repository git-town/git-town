# Creates and pushes a branch
def create_branch branch_name
  run "git checkout -b #{branch_name} main"
  run "git push -u origin #{branch_name}"
end


# Returns the name of the branch that is currently checked out
def current_branch_name
  output_of("git rev-parse --abbrev-ref HEAD")
end


# Returns the names of the existing feature branches
def existing_feature_branches
  existing_local_branches - %w[main]
end

# Returns the names of all existing local branches.
#
# Does not return the "master" branch nor remote branches.
#
# The branches are ordered this ways:
# * main branch
# * feature branches ordered alphabetically
def existing_local_branches
  actual_branches = array_output_of("git branch | tr -d '*'")
  actual_branches.delete('master')
  actual_main_branch = actual_branches.delete 'main'
  [actual_main_branch].concat(actual_branches)
end


# Returns the names of all existing remote branches.
#
# Does not return the "master" branch.
def existing_remote_branches
  remote_branches = run('git branch -a | grep remotes').out
                                                       .split("\n")
                                                       .map(&:strip)
  remote_branches.delete('remotes/origin/master')
  remote_branches.delete('remotes/origin/HEAD -> origin/master')
  remote_branches
end

def number_of_branches_out_of_sync
  integer_output_of("git branch -vv | grep -o '\[.*\]' | tr -d '[]' | awk '{ print $2 }' | grep . | wc -l")
end


# Executes the given block, then returns to the currently checked out branch
def returning_to_current_branch
  original_branch = current_branch_name
  yield
  run "git checkout #{original_branch}"
end


# Verifies the branches in each repository
def verify_branches branches_array
  branches_array.each do |branches|
    branch_names = Kappamaki.from_sentence(branches['branches'])

    case branches['repository']
    when 'local'
      expect(existing_local_branches).to match_array(branch_names)
    when 'remote'
      branch_names = branch_names.map { |n| "remotes/origin/#{n}" }
      expect(existing_remote_branches).to match_array(branch_names)
    when 'coworker'
      at_path coworker_repository_path do
        expect(existing_local_branches).to match_array(branch_names)
      end
    else
      raise "Unknown repository: #{branches['repository']}"
    end
  end
end
