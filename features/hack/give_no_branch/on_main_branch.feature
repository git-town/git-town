Feature: cannot make the current main branch a feature branch

  Background:
    Given a Git repo with origin
    When I run "git-town hack"

  Scenario: result
    Then Git Town runs no commands
    And it prints the error:
      """
      you are trying to convert the main branch to a feature branch. That's not possible. If you want to create a feature branch, did you forget to add the branch name?
      """
    And the main branch is still "main"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the main branch is still "main"
