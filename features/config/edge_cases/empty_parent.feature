Feature: ignore empty parent branch setting

  Scenario:
    And Git Town parent setting for branch "foo" is ""
    When I run "git-town config"
    Then it prints:
      """
      NOTICE: I have found an empty parent configuration entry for branch "foo".
      I have deleted this configuration entry.
      """
