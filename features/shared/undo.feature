Feature: cannot double undo

  As a developer accidently running undo twice in a row
  I want to be warned that there is nothing to undo
  So that it does undo something twice (most likely causing errors) or undo the undo

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
