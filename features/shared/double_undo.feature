Feature: no double undo

  Scenario:
    Given my repo has a feature branch "feature"
    And I am on the "feature" branch
    And I run "git-town kill"
    And I run "git-town undo"
    When I run "git-town undo"
    Then it prints the error:
      """
      nothing to undo
      """
