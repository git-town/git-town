Feature: show all available commands

  Scenario: run without a command
    When I run "git-town"
    Then it prints:
      """
      Available Commands:
      """
