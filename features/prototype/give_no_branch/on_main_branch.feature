Feature: prototype the current main branch

  Background:
    Given a Git repo with origin
    When I run "git-town prototype"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      cannot prototype the main branch
      """
    And the current branch is still "main"
    And the main branch is still "main"
    And there are still no prototype branches

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "main"
    And the main branch is still "main"
    And there are still no prototype branches
