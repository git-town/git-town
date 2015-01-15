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

  Scenario Outline: Running outside of a git repository
    Given I'm currently not in a git repository
    When I run `<COMMAND>`
    Then I see the "git-town" man page
    And I don't see "fatal: Not a git repository"

    Examples:
      | COMMAND       |
      | git town      |
      | git town help |
