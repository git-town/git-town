def delete_configuration
  run_this 'git config --unset git-town.main-branch-name'
end
