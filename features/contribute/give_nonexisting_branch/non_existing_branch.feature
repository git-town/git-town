Feature: cannot make non-existing branches contribution branches

  Background:
    Given a Git repo with origin
    When I run "git-town contribute non-existing"

  Scenario: result
    Then Git Town runs no commands
    And it prints the error:
      """
      there is no branch "non-existing"
      """
    And the current branch is still "main"
    And there are still no contribution branches

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And there are still no contribution branches
    And the current branch is still "main"
