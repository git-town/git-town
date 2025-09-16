Feature: detaching in headless state

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | branch-1 | feature | main   | local     |
    And the commits
      | BRANCH   | LOCATION | MESSAGE   |
      | branch-1 | local    | commit 1a |
      | branch-1 | local    | commit 1b |
    And the current branch is "branch-1"
    And I run "git checkout HEAD^"
    When I run "git-town detach"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      cannot determine current branch
      """
    And the initial lineage exists now
