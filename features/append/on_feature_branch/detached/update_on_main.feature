Feature: append a new feature branch to an existing feature branch in detached mode when there are updates on the main branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | existing | local, origin | existing commit |
      | main     | local, origin | main commit     |
    And the current branch is "existing"
    When I run "git-town append new --detached"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                  |
      | existing | git fetch --prune --tags                 |
      |          | git merge --no-edit --ff main            |
      |          | git merge --no-edit --ff origin/existing |
      |          | git push                                 |
      |          | git checkout -b new                      |
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE                           |
      | main     | local, origin | main commit                       |
      | existing | local, origin | existing commit                   |
      |          |               | Merge branch 'main' into existing |
    And this lineage exists now
      | BRANCH   | PARENT   |
      | existing | main     |
      | new      | existing |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | new      | git checkout existing                           |
      | existing | git reset --hard {{ sha 'existing commit' }}    |
      |          | git push --force-with-lease --force-if-includes |
      |          | git branch -D new                               |
    And the initial commits exist now
    And the initial lineage exists now
