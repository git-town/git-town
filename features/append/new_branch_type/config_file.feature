Feature: append a new branch when prototype branches are configured via the config file

  Background:
    Given a Git repo with origin
    And the committed configuration file:
      """
      [create]
      new-branch-type = "prototype"
      """
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | existing | local, origin | existing commit |
    And the current branch is "existing"
    When I run "git-town append new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                  |
      | existing | git fetch --prune --tags |
      |          | git checkout -b new      |
    And this lineage exists now
      """
      main
        existing
          new
      """
    And branch "new" now has type "prototype"
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND               |
      | new      | git checkout existing |
      | existing | git branch -D new     |
    And the initial lineage exists now
    And the initial commits exist now
