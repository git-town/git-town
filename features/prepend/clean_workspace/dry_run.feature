Feature: dry-run prepending a branch to a feature branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the current branch is "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, origin | old commit |
    When I run "git-town prepend parent --dry-run"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                             |
      | old    | git fetch --prune --tags            |
      |        | git checkout main                   |
      | main   | git rebase origin/main              |
      |        | git checkout old                    |
      | old    | git merge --no-edit --ff origin/old |
      |        | git merge --no-edit --ff main       |
      |        | git checkout -b parent main         |
    And the current branch is still "old"
    And the initial commits exist
    And the initial branches and lineage exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "old"
    And the initial commits exist
    And the initial lineage exists
