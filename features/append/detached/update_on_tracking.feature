Feature: append a new feature branch to an existing feature branch in detached mode when there are updates to its tracking branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | existing | local, origin | existing commit |
      | existing | origin        | origin commit   |
    And the current branch is "existing"
    When I run "git-town append new --detached"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                  |
      | existing | git fetch --prune --tags                 |
      |          | git merge --no-edit --ff origin/existing |
      |          | git checkout -b new                      |
    And this lineage exists now
      """
      main
        existing
          new
      """
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE         |
      | existing | local, origin | existing commit |
      |          |               | origin commit   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                      |
      | new      | git checkout existing                        |
      | existing | git reset --hard {{ sha 'existing commit' }} |
      |          | git branch -D new                            |
    And the initial lineage exists now
    And the initial commits exist now
