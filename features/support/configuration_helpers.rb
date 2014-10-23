def delete_main_branch_configuration
  run 'git config --unset git-town.main-branch-name'
end

def set_non_feature_branches_configuration configuration
  run "git config git-town.non-feature-branch-names #{configuration}"
end
