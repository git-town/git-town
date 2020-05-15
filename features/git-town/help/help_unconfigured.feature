Feature: show help screen when Git Town is not configured

  (see ./help_configured.feature)

  Scenario Outline:
    Given I haven't configured Git Town yet
    When I run "<COMMAND>"
    Then it prints
      """
      Usage:
        git-town [command]
      """

    Examples:
      | COMMAND       |
      | git-town      |
      | git-town help |
