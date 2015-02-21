def configure_non_feature_branches configuration
  set_configuration 'non-feature-branch-names', configuration
end


def delete_main_branch_configuration
  run 'git config --unset git-town.main-branch-name'
end


def delete_non_feature_branches_configuration
  run 'git config --unset git-town.non-feature-branch-names'
end


def git_town_configuration
  # OR'ed with true so that this doesn't exit with an error if config doesn't exist
  array_output_of 'git config --get-regex git-town || true'
end


def main_branch_configuration
  output_of 'git config --get git-town.main-branch-name || true'
end


def non_feature_branch_configuration
  output_of 'git config --get git-town.non-feature-branch-names || true'
end


def set_configuration configuration, value
  run "git config git-town.#{configuration} '#{value}'"
end
