Feature: cannot make non-existing branches contribution branches

  Background:
    Given a Git repo with origin
    When I run "git-town contribute non-existing"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      there is no branch "non-existing"
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
