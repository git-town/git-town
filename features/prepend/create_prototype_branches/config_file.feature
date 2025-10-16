Feature: prepend a new branch when prototype branches are configured via config file

  Background:
    Given a Git repo with origin
    And the committed configuration file:
      """
      [create]
      new-branch-type = "prototype"
      """
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, origin | old commit |
    And the current branch is "old"
    When I run "git-town prepend parent"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                     |
      | old    | git fetch --prune --tags    |
      |        | git checkout -b parent main |
    And this lineage exists now
      """
      main
        parent
          old
      """
    And the initial commits exist now
    And branch "parent" now has type "prototype"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND              |
      | parent | git checkout old     |
      | old    | git branch -D parent |
    And the initial lineage exists now
    And the initial commits exist now
