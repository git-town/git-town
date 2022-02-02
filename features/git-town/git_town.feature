Feature: show all available commands

  Scenario: no command given
    When I run "git-town"
    Then it prints:
      """
      Available Commands:
      """
