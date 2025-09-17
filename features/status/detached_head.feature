Feature: describe the status of a detached head

  Background:
    Given a Git repo with origin
    And I ran "git-town sync"
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS |
      | branch | feature | main   | local     |
    And the commits
      | BRANCH | LOCATION | MESSAGE  |
      | branch | local    | commit 1 |
      |        | local    | commit 2 |
    And the current branch is "branch"
    And I ran "git checkout HEAD^"
    When I run "git-town status"

  Scenario: result
    Then Git Town prints:
      """
      The previous Git Town command (sync) finished successfully.
      You can run "git town undo" to go back to where you started.
      """
