Feature: cannot move a detached head up

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | branch | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION | MESSAGE  |
      | branch | local    | commit 1 |
      | branch | local    | commit 2 |
    And the current branch is "branch"
    And I ran "git checkout HEAD^"
    When I run "git-town up"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND |
    And Git Town prints the error:
      """
      you need to be on a branch to go up
      """
