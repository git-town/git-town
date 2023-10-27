Feature: empty parent branch setting

  @debug @this
  Scenario:
    Given the current branch is a feature branch "feature"
    And local setting "git-town-branch.branch-name.parent" is ""
    When I run "git-town sync"
    Then it prints
      """
      xxx
      """
