Feature: Show correct git town usage

  Scenario: invalid git town command
    When I run `git-town invalidcommand`
    Then I get the error
      """
      Error: unknown command "invalidcommand" for "git-town"
      Run 'git-town --help' for usage.
      """
