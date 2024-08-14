Feature: dry-run deleting the current feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | current | feature | main   | local, origin |
      | other   | feature | main   | local, origin |
    And the current branch is "current"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | current | local, origin | current commit |
      | other   | local, origin | other commit   |
    And an uncommitted file
    And the current branch is "current" and the previous branch is "other"
    When I run "git-town kill --dry-run"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                                          |
      | current | git fetch --prune --tags                         |
      |         | git push origin :current                         |
      |         | git add -A                                       |
      |         | git commit -m "Committing WIP for git town undo" |
      |         | git checkout other                               |
      | other   | git branch -D current                            |
    And the current branch is still "current"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial branches and lineage exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is now "current"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial branches and lineage exist
