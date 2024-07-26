Feature: cannot make the current main branch a feature branch

  Background:
    Given a Git repo with origin
    When I run "git-town hack"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      cannot make the main branch a feature branch
      """
    And the main branch is still "main"

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the main branch is still "main"
