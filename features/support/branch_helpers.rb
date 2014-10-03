# Returns the name of the branch that is currently checked out
def current_branch_name
  run("git rev-parse --abbrev-ref HEAD")[:out]
end


# Returns the names of all existing local branches.
#
# Does not return the "master" branch nor remote branches.
#
# The branches are ordered this ways:
# * main branch
# * feature branches ordered alphabetically
def existing_local_branches
  actual_branches = run("git branch").fetch(:out)
                                          .split("\n")
                                          .map(&:strip)
                                          .map{|s| s.sub('* ', '')}
  actual_branches.delete('master')
  actual_main_branch = actual_branches.delete 'main'
  [actual_main_branch].concat(actual_branches)
                      .compact
end


# Returns the names of all existing remote branches.
#
# Does not return the "master" branch.
#
# The branches are ordered this ways:
# * main branch
# * feature branches ordered alphabetically
def existing_remote_branches
  remote_branches = run('git branch -a | grep remotes').fetch(:out)
                                                       .split("\n")
                                                       .map(&:strip)
  remote_branches.delete('remotes/origin/master')
  remote_branches.delete('remotes/origin/HEAD -> origin/master')
  remote_branches
end


def remote_branch_exists branch_name
  run("git branch -a | grep remotes/origin/#{branch_name} | wc -l")[:out] != '0'
end
