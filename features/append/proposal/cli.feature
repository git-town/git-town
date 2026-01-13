Feature: append a new feature branch to an existing feature branch and sync proposals

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | existing | local, origin | existing commit |
    And the proposals
      | ID | SOURCE BRANCH | TARGET BRANCH | BODY |
      | 1  | existing      | main          | body |
    And Git setting "git-town.proposals-show-lineage" is "cli"
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
    And the initial commits exist now
    And the initial proposals exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND               |
      | new      | git checkout existing |
      | existing | git branch -D new     |
    And the initial lineage exists now
    And the initial commits exist now
    And the initial proposals exist now
