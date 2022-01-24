Feature: show help screen when Git Town is configured

  Scenario Outline:
    When I run "<COMMAND>"
    Then it prints:
      """
      Usage:
        git-town [command]
      """

    Examples:
      | COMMAND       |
      | git-town      |
      | git-town help |
