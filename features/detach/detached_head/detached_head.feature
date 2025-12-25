Feature: cannot detach a detached head

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | branch-1 | feature | main   | local     |
    And the commits
      | BRANCH   | LOCATION | MESSAGE  |
      | branch-1 | local    | commit 1 |
      | branch-1 | local    | commit 2 |
    And the current branch is "branch-1"
    And I run "git checkout HEAD^"
    When I run "git-town detach"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH                     | COMMAND                  |
      | {{ sha-short 'commit 1' }} | git fetch --prune --tags |
    And Git Town prints the error:
      """
      please check out the branch to detach
      """
    And the initial lineage exists now
