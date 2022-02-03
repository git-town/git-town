Feature: does not double undo

  Scenario: call undo twice
    Given my repo has a feature branch "feature"
    And I am on the "feature" branch
    And I run "git-town kill"
    And I am now on the "main" branch
    And I run "git-town undo"
    And I am now on the "feature" branch
    When I run "git-town undo"
    Then it prints the error:
      """
      nothing to undo
      """
