Feature: show all available commands

  Background:
    Given I am outside a Git repo

  Scenario: no command
    When I run "git-town"
    Then Git Town prints:
      """
      Basic commands:
      """

  Scenario: unknown command
    When I run "git-town invalidcommand"
    Then Git Town prints the error:
      """
      Error: unknown command "invalidcommand" for "git-town"
      """
