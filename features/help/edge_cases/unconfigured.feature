Feature: show help even if the current repo misses configuration

  Scenario Outline:
    Given I am outside a Git repo
    When I run "<COMMAND>"
    Then it prints:
      """
      Usage:
        git-town [flags]
        git-town [command]
      """

    Examples:
      | COMMAND       |
      | git-town      |
      | git-town help |
