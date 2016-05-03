def configure_perennial_branches configuration
  set_configuration 'perennial-branch-names', configuration
end


def delete_main_branch_configuration
  run 'git config --unset git-town.main-branch-name'
end


def delete_perennial_branches_configuration
  run 'git config --unset git-town.perennial-branch-names'
end


def get_configuration configuration
  output_of "git config --get git-town.#{configuration} || true"
end


def git_town_configuration
  # OR'ed with true so that this doesn't exit with an error if config doesn't exist
  array_output_of 'git config --get-regex git-town || true'
end


def main_branch_configuration
  get_configuration 'main-branch-name'
end


def non_feature_branch_configuration
  get_configuration 'non-feature-branch-names'
end


def perennial_branch_configuration
  get_configuration 'perennial-branch-names'
end


def pull_branch_strategy_configuration
  get_configuration 'pull-branch-strategy'
end


def set_configuration configuration, value
  run "git config git-town.#{configuration} '#{value}'"
end
