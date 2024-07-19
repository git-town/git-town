Feature: remove parent entries for perennial branches

  Scenario: child is its own parent
    Given a Git repo clone
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And Git Town parent setting for branch "feature" is "feature"
    When I run "git town config"
    Then it prints:
      """
      removing lineage entry for "feature" because the parent is the child
      """
    And no lineage exists now
