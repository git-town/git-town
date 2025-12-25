Feature: delete the given branch from a detached head

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
      | branch-2 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION | MESSAGE   |
      | branch-1 | local    | commit 1a |
      | branch-1 | local    | commit 1b |
      | branch-2 | local    | commit 2a |
    And the current branch is "branch-1"
    And I ran "git checkout HEAD^"
    When I run "git-town delete branch-2"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH                      | COMMAND                  |
      | {{ sha-short 'commit 1a' }} | git fetch --prune --tags |
    And Git Town prints the error:
      """
      please check out the branch to delete
      """
    And the initial branches and lineage exist now
    And the initial commits exist now
