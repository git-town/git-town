Feature: rename a branch to itself

  Background:
    Given my repo has a feature branch "feature"
    And I am on the "feature" branch

  Scenario: without force
    When I run "git-town rename-branch feature"
    Then it runs no commands
    And it prints the error:
      """
      cannot rename branch to current name
      """

  Scenario: with force
    When I run "git-town rename-branch --force feature"
    Then it runs no commands
    And it prints the error:
      """
      cannot rename branch to current name
      """
