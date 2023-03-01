Feature: show all available commands

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
      Run 'git-town --help' for usage.
      """
