Feature: Show correct git town usage

  Scenario: invalid git town command
    When I run `gt invalidcommand`
    Then I get the error
      """
      Error: unknown command "invalidcommand" for "gt"
      Run 'gt --help' for usage.
      """
