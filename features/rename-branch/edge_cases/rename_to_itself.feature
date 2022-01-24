Feature: rename a branch to itself

  Scenario: without force
    Given my repo has a feature branch named "feature"
    Given I am on the "feature" branch
    When I run "git-town rename-branch feature"
    Then it runs no commands
    And it prints the error:
      """
      cannot rename branch to current name
      """

  Scenario: with force
    Given my repo has a feature branch named "feature"
    Given I am on the "feature" branch
    When I run "git-town rename-branch --force feature"
    Then it runs no commands
    And it prints the error:
      """
      cannot rename branch to current name
      """
