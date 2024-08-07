@smoke
Feature: don't sync tags while appending

  Background:
    Given a Git repo with origin
    And the branch
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | existing | local, origin | existing commit |
    And the current branch is "existing"
    And Git Town setting "sync-tags" is "false"
    When I run "git-town append new"

  @this
  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                                  |
      | existing | git fetch --prune --no-tags              |
      |          | git checkout main                        |
      | main     | git rebase origin/main                   |
      |          | git checkout existing                    |
      | existing | git merge --no-edit --ff origin/existing |
      |          | git merge --no-edit --ff main            |
      |          | git checkout -b new                      |
    And the current branch is now "new"
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE         |
      | existing | local, origin | existing commit |
      | new      | local         | existing commit |
    And this lineage exists now
      | BRANCH   | PARENT   |
      | existing | main     |
      | new      | existing |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH   | COMMAND               |
      | new      | git checkout existing |
      | existing | git branch -D new     |
    And the current branch is now "existing"
    And the initial commits exist
    And the initial lineage exists
