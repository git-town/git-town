Feature: park the current detached state

  Background:
    Given a local Git repo
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS |
      | branch | feature | main   | local     |
    And the commits
      | BRANCH | LOCATION | MESSAGE  |
      | branch | local    | commit 1 |
      |        | local    | commit 2 |
    And the current branch is "branch"
    And I ran "git checkout HEAD^"
    When I run "git-town park"
  # TODO: fix the broken behavior: it should not park here

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints something like:
      """
      branch .* is now parked
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "branch" still has type "feature"
