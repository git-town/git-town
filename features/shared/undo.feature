Feature: cannot double undo


  Background:
    Given my repo has a feature branch named "feature"
    And I am on the "feature" branch
    And I run "git-town kill"
    And I am now on the "main" branch
    And I run "git-town undo"
    And I am now on the "feature" branch

  Scenario:
    When I run "git-town undo"
    Then it prints the error:
      """
      nothing to undo
      """
