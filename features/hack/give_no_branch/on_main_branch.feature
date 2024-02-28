Feature: cannot make the current main branch a feature branch

  Background:
    Given the current branch is "main"
    When I run "git-town hack"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      branch "main" is the main branch
      """
    And the main branch is still "main"

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the main branch is still "main"
