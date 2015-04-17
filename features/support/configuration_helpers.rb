def configure_git user, secondary: false
  run "git config user.name #{user}"
  run "git config user.email #{user}@#{secondary ? 'secondary' : 'primary'}.com"
  run 'git config push.default simple'
  run 'git config core.editor vim'

  # Git Town Configuration
  run 'git config git-town.main-branch-name main'
  run 'git config git-town.non-feature-branch-names ""'
end
