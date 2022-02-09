Feature: no double undo

  Scenario:
    Given a feature branch "feature"
    And the current branch is "feature"
    And I run "git-town kill"
    And I run "git-town undo"
    When I run "git-town undo"
    Then it prints the error:
      """
      nothing to undo
      """
