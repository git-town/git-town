Feature: ship the given branch from a detached head

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | branch-1 | feature | main   | local     |
      | branch-2 | feature | main   | local     |
    And the commits
      | BRANCH   | LOCATION | MESSAGE   |
      | branch-1 | local    | commit 1a |
      |          | local    | commit 1b |
      | branch-2 | local    | commit 2  |
    And the current branch is "branch-1"
    And I ran "git checkout HEAD^"
    When I run "git-town ship branch-2"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH                      | COMMAND                  |
      | {{ sha-short 'commit 1a' }} | git fetch --prune --tags |
    And Git Town prints the error:
      """
      please check out the branch to ship
      """
