Feature: inside a committed subfolder that exists only on the current feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE       | FILE NAME              |
      | existing | local, origin | folder commit | committed_folder/file1 |
    And the current branch is "existing"
    When I run "git-town hack new" in the "committed_folder" folder

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                  |
      | existing | git fetch --prune --tags |
      |          | git checkout main        |
      | main     | git checkout -b new      |
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
    And the initial commits exist now
