Feature: detaching a branch that is behind its tracking branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE   |
      | branch-1 | local, origin | commit 1a |
      | branch-1 | local, origin | commit 1b |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE   |
      | branch-2 | local, origin | commit 2a |
      | branch-2 | origin        | commit 2b |
    And the current branch is "branch-2"
    When I run "git-town detach"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                  |
      | branch-2 | git fetch --prune --tags |
    And Git Town prints the error:
      """
      please sync this stack before detaching branches from it
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "branch-2"
    And the initial commits exist now
    And the initial lineage exists now
