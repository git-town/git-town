Feature: cannot park non-existing branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the current branch is "feature"
    When I run "git-town park feature non-existing"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      there is no branch "non-existing"
      """
    And branch "feature" still has type "feature"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "feature" still has type "feature"
