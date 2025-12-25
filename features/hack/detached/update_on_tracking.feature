Feature: on a feature branch with a clean workspace in detached mode with updates on the tracking branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION | MESSAGE                |
      | existing | local    | existing local commit  |
      | existing | origin   | existing origin commit |
    And the current branch is "existing"
    When I run "git-town hack new --detached"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                  |
      | existing | git fetch --prune --tags |
      |          | git checkout -b new main |
    And this lineage exists now
      """
      main
        existing
        new
      """
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND               |
      | new      | git checkout existing |
      | existing | git branch -D new     |
    And the initial branches and lineage exist now
    And the initial commits exist now
