Feature: describe

  Scenario: no runstate
    When I run "git-town status"
    Then it prints:
      """
      No status file found for this repository.
      """
