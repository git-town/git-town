Feature: cannot observe non-existing branches

  Background:
    Given a Git repo with origin
    When I run "git-town observe non-existing"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      there is no branch "non-existing"
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
