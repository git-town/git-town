Feature: Help for unknown command

  Scenario: invalid git town command
    When I run "git-town invalidcommand"
    Then it prints the error:
      """
      Error: unknown command "invalidcommand" for "git-town"
      Run 'git-town --help' for usage.
      """
