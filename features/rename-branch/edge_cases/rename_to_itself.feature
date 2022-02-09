Feature: rename a branch to itself

  Background:
    Given a feature branch "old"
    And the current branch is "old"

  Scenario: without force
    When I run "git-town rename-branch old"
    Then it runs no commands
    And it prints the error:
      """
      cannot rename branch to current name
      """

  Scenario: with force
    When I run "git-town rename-branch --force old"
    Then it runs no commands
    And it prints the error:
      """
      cannot rename branch to current name
      """
