Feature: must provide a commit message

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the current branch is "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | local    | feature commit |
    And Git Town setting "ship-strategy" is "squash-merge"
    When I run "git-town ship" and close the editor

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                         |
      | feature | git fetch --prune --tags        |
      |         | git checkout main               |
      | main    | git merge --squash --ff feature |
      |         | git commit                      |
      |         | git reset --hard                |
      |         | git checkout feature            |
    And it prints the error:
      """
      aborted because merge exited with error
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
