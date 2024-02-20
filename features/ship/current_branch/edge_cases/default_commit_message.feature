Feature: must provide a commit message

  Background:
    Given the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | local    | feature commit |
    When I run "git-town ship" and close the editor

  @this
  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                    |
      | feature | git fetch --prune --tags   |
      |         | git checkout main          |
      | main    | git merge --squash feature |
      |         | git commit                 |
      |         | git reset --hard           |
      |         | git checkout feature       |
    And it prints the error:
      """
      aborted because commit exited with error
      """
    And the current branch is still "feature"
    And the initial commits exist
    And the initial lineage exists

  Scenario: undo
    When I run "git-town undo"
    Then it prints:
      """
      nothing to undo
      """
    And the current branch is still "feature"
    And the initial commits exist
    And the initial lineage exists
