Feature: Show correct git town usage

  Scenario: invalid git town command
    When I run `git town invalidcommand`
    Then I see
      """
      error: unsupported subcommand 'invalidcommand'
      usage: git town
         or: git town config
         or: git town help
         or: git town main-branch [<branchname>]
         or: git town non-feature-branches [(--add | --remove) <branchname>]
         or: git town version
      """
