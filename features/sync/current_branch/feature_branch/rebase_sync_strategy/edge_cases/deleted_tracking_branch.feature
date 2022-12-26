Feature: restores deleted tracking branch

  Background:
    Given setting "sync-strategy" is "rebase"
    And the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And origin deletes the "feature" branch
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                    |
      | feature | git fetch --prune --tags   |
      |         | git checkout main          |
      | main    | git rebase origin/main     |
      |         | git checkout feature       |
      | feature | git rebase main            |
      |         | git push -u origin feature |
    And all branches are now synchronized
    And the current branch is still "feature"
    And now the initial commits exist
