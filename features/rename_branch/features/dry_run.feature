Feature: rename the current branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the current branch is "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | old    | local, origin | old commit  |
    When I run "git-town rename-branch new --dry-run"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | old    | git fetch --prune --tags |
      |        | git branch new old       |
      |        | git checkout new         |
      | new    | git push -u origin new   |
      |        | git push origin :old     |
      |        | git branch -D old        |
    And the current branch is still "old"
    And the initial commits exist
    And the initial branches and lineage exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "old"
    And the initial commits exist
    And the initial branches and lineage exist
