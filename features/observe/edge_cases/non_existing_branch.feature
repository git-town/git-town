Feature: cannot observe non-existing branches

  Background:
    When I run "git-town observe non-existing"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      remote branch "origin/non-existing" doesn't seem to exist
      """
    And the current branch is still "main"
    And there are still no observed branches

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And there are still no observed branches
    And the current branch is still "main"
