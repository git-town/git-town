def delete_configuration
  run 'git config --unset git-town.main-branch-name'
end
