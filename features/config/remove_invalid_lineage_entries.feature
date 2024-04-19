Feature: remove parent entries for perennial branches

  Scenario: child is its own parent
    Given the local feature branch "feature"
    And Git Town parent setting for branch "feature" is "feature"
    When I run "git town config"
    Then it prints:
      """
      removing lineage entry for "feature" because the parent is the child
      """
    And no lineage exists now
      | BRANCH    | PARENT |
      | feature-2 | main   |
