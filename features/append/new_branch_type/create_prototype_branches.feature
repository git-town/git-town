Feature: append a new branch when prototype branches are configured via a deprecated setting in Git metadata

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | existing | local, origin | existing commit |
    And Git setting "git-town.create-prototype-branches" is "true"
    And the current branch is "existing"
    When I run "git-town append new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                  |
      | existing | git fetch --prune --tags |
      |          | git checkout -b new      |
    And Git Town prints:
      """
      Upgrading deprecated local setting git-town.create-prototype-branches to git-town.new-branch-type
      """
    And Git setting "git-town.create-prototype-branches" now doesn't exist
    And Git setting "git-town.new-branch-type" is now "prototype"
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
    And Git setting "git-town.create-prototype-branches" still doesn't exist
    And Git setting "git-town.new-branch-type" is still "prototype"
    And the initial lineage exists now
    And the initial commits exist now
