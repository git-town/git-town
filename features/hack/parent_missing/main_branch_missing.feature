Feature: on a feature branch but without main branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION | MESSAGE         |
      | main     | origin   | main commit     |
      | existing | local    | existing commit |
    And the current branch is "existing"
    And I ran "git branch -d main"
    When I run "git-town hack new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                    |
      | existing | git fetch --prune --tags                   |
      |          | git checkout -b new origin/main --no-track |
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
    And the initial lineage exists now
    And the branches are now
      | REPOSITORY | BRANCHES       |
      | local      | existing       |
      | origin     | main, existing |
    And the initial commits exist now
