Feature: no double undo

  Scenario:
    Given the current branch is a feature branch "feature"
    And I run "git-town kill"
    And I run "git-town undo"
    When I run "git-town undo"
    Then it prints:
      """
      nothing to undo
      """
