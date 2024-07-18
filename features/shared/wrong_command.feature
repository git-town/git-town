Feature: show all available commands

  Background:
    Given I am outside a Git repo

  Scenario: no command
    When I run "git-town"
    Then it prints:
      """
      Basic commands:
      """

  Scenario: unknown command
    When I run "git-town invalidcommand"
    Then it prints the error:
      """
      Error: unknown command "invalidcommand" for "git-town"
      """
