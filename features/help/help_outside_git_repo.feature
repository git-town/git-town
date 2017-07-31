Feature: show help screen even outside of a Git repository

  (see ./help_configured.feature)


  Scenario Outline: Running outside of a Git repository
    Given I'm currently not in a git repository
    When I run `<COMMAND>`
    Then I see
      """
      Usage:
        git-town [command]
      """
    And I don't see "fatal: Not a Git repository"

    Examples:
      | COMMAND       |
      | git-town      |
      | git-town help |
