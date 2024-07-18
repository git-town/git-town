Feature: cannot park non-existing branches

  Background:
    Given a Git repo clone
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the current branch is "feature"
    And an uncommitted file
    When I run "git-town park feature non-existing"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      there is no branch "non-existing"
      """
    And the current branch is still "feature"
    And the uncommitted file still exists
    And there are still no parked branches

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And there are still no parked branches
    And the current branch is still "feature"
