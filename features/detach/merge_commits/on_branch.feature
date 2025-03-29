Feature: detaching a branch that contains merge commits

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-1 | feature | main     | local, origin |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  |
      | branch-1 | local, origin | commit 1 |
      | branch-2 | local, origin | commit 2 |
    And the current branch is "branch-2"
    And I ran "git merge branch-1 --no-edit"
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-3 | feature | branch-2 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  |
      | branch-3 | local, origin | commit 3 |
    And the current branch is "branch-2"
    When I run "git-town detach"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                  |
      | branch-2 | git fetch --prune --tags |
    And Git Town prints the error:
      """
      branch "branch-2" contains merge commits, please compress and try again
      """
    And the current branch is still "branch-2"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "branch-2"
