Feature: Show correct git town usage

  Scenario: invalid git town command
    When I run `git town invalidcommand`
    Then I get the error
      """
      error: unsupported subcommand 'invalidcommand'
      usage: git town
         or: git town config [--reset | --setup]
         or: git town hack-push-flag [(true | false)]
         or: git town help
         or: git town install-fish-autocompletion
         or: git town main-branch [<branch_name>]
         or: git town perennial-branches [(--add | --remove) <branch_name>]
         or: git town pull-branch-strategy [(rebase | merge)]
         or: git town version
      """
