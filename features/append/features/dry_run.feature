Feature: dry run appending a new feature branch to an existing feature branch

  Background:
    Given the current branch is a feature branch "existing"
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | existing | local, origin | existing commit |
    And an uncommitted file
    When I run "git-town append new --dry-run"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                             |
      | existing | git fetch --prune --tags            |
      |          | git add -A                          |
      |          | git stash                           |
      |          | git checkout main                   |
      | main     | git rebase origin/main              |
      |          | git checkout existing               |
      | existing | git merge --no-edit origin/existing |
      |          | git merge --no-edit main            |
      |          | git branch new existing             |
      |          | git checkout new                    |
      | new      | git stash pop                       |
    And the current branch is still "existing"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial branches and lineage exist

  @debug @this
  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "existing"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial lineage exists
