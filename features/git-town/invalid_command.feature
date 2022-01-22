Feature: Help for unknown command

  To learn to use Git Town correctly
  When running an unknown command
  I want to see guidance telling me that this is an unknown command.

  Scenario: invalid git town command
    When I run "git-town invalidcommand"
    Then it prints the error:
      """
      Error: unknown command "invalidcommand" for "git-town"
      Run 'git-town --help' for usage.
      """
