Feature: show help screen even outside of a Git repository

  (see ./help_configured.feature)

  Scenario Outline: Running outside of a Git repository
    Given my workspace is currently not a Git repo
    When I run "<COMMAND>"
    Then it prints:
      """
      Usage:
        git-town [command]
      """
    And it does not print "fatal: Not a Git repository"

    Examples:
      | COMMAND       |
      | git-town      |
      | git-town help |
