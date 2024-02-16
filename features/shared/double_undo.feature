Feature: no double undo

  @this

  Scenario:
    Given the current branch is a feature branch "feature"
    And I run "git-town kill"
    And I run "git-town undo"
    When I run "git-town undo"
    Then it prints:
      """
      nothing to undo
      """
