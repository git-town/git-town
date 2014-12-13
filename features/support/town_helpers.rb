def configure_non_feature_branches configuration
  set_configuration 'non-feature-branch-names', configuration
end


def delete_main_branch_configuration
  run 'git config --unset git-town.main-branch-name'
end


def delete_non_feature_branches_configuration
  run 'git config --unset git-town.non-feature-branch-names'
end


def main_branch_configuration
  output_of 'git config --get git-town.main-branch-name'
end


def non_feature_branch_configuration
  output_of 'git config --get git-town.non-feature-branch-names'
end


def set_configuration configuration, value
  run "git config git-town.#{configuration} '#{value}'"
end
