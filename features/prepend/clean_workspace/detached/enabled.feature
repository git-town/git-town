Feature: prepend a branch to a feature branch in detached mode

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the current branch is "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, origin | old commit |
    When I run "git-town prepend parent --detached"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                             |
      | old    | git fetch --prune --tags            |
      |        | git merge --no-edit --ff main       |
      |        | git merge --no-edit --ff origin/old |
      |        | git checkout -b parent main         |
    And the current branch is now "parent"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, origin | old commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | old    | parent |
      | parent | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND              |
      | parent | git checkout old     |
      | old    | git branch -D parent |
    And the current branch is now "old"
    And the initial commits exist now
    And the initial lineage exists now
