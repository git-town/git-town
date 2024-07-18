Feature: request help for an unknown command

  Scenario: unknown command
    Given I am outside a Git repo
    When I run "git-town help zonk"
    Then it prints:
      """
      Unknown help topic
      """
