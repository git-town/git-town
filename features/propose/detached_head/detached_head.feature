Feature: Cannot create proposals in detached mode

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS |
      | branch | feature | main   | local     |
    And the commits
      | BRANCH | LOCATION | MESSAGE  |
      | branch | local    | commit 1 |
      |        | local    | commit 2 |
    And the current branch is "branch"
    And I ran "git checkout HEAD^"
    When I run "git-town propose"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      cannot propose in detached head state
      """
