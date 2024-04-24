Feature: dry run appending a new feature branch to an existing feature branch

  Background:
    Given the current branch is a feature branch "existing"
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | existing | local, origin | existing commit |
    When I run "git-town append new --dry-run"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                                  |
      | existing | git fetch --prune --tags                 |
      |          | git checkout main                        |
      | main     | git rebase origin/main                   |
      |          | git checkout existing                    |
      | existing | git merge --no-edit --ff origin/existing |
      |          | git merge --no-edit --ff main            |
      |          | git checkout -b new                      |
    And the current branch is still "existing"
    And the initial commits exist
    And the initial branches and lineage exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "existing"
    And the initial commits exist
    And the initial lineage exists
