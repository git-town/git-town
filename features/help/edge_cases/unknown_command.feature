Feature: request help for an unknown command

  Scenario: unknown command
    When I run "git-town help zonk"
    Then it prints:
      """
      Unknown help topic
      """
