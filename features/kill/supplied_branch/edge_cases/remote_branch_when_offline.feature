Feature: does not kill a remote branch in offline mode

  Background:
    Given a remote feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | origin   | feature commit |
    And I run "git fetch"
    And offline mode is enabled
    And the current branch is "main"
    When I run "git-town kill feature"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      cannot delete remote branch "origin/feature" in offline mode
      """
    And the current branch is still "main"
    And no branch hierarchy exists now
