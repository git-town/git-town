Feature: prepend a new branch when prototype branches are configured via config file

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the current branch is "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, origin | old commit |
    And the committed configuration file:
      """
      new-branch-type = "prototype"
      """
    When I run "git-town prepend parent"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                 |
      | old    | git fetch --prune --tags                |
      |        | git checkout main                       |
      | main   | git rebase origin/main --no-update-refs |
      |        | git checkout old                        |
      | old    | git merge --no-edit --ff main           |
      |        | git merge --no-edit --ff origin/old     |
      |        | git checkout -b parent main             |
    And the current branch is now "parent"
    And branch "parent" is now prototype
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, origin | old commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | old    | parent |
      | parent | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND              |
      | parent | git checkout old     |
      | old    | git branch -D parent |
    And the current branch is now "old"
    And the initial commits exist now
    And the initial lineage exists now
