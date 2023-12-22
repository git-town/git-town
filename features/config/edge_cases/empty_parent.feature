Feature: ignore empty parent branch setting

  Scenario:
    And Git Town parent setting for branch "foo" is ""
    When I run "git-town config"
    Then it prints:
      """
      NOTICE: deleted empty configuration entry "git-town-branch.foo.parent"
      """
