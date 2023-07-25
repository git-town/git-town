Feature: auto-push new branches

  Background:
    Given setting "push-new-branches" is "true"
    And the current branch is a feature branch "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE        |
      | old    | local, origin | feature commit |
    When I run "git-town prepend new"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                        |
      | old    | git fetch --prune --tags       |
      |        | git checkout main              |
      | main   | git rebase origin/main         |
      |        | git checkout old               |
      | old    | git merge --no-edit origin/old |
      |        | git merge --no-edit main       |
      |        | git branch new main            |
      |        | git checkout new               |
      | new    | git push -u origin new         |
    And the current branch is now "new"
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE        |
      | old    | local, origin | feature commit |
    And this branch lineage exists now
      | BRANCH | PARENT |
      | new    | main   |
      | old    | new    |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND              |
      | new    | git push origin :new |
      |        | git checkout old     |
      | old    | git branch -D new    |
      |        | git checkout main    |
      | main   | git checkout old     |
    And the current branch is now "old"
    And now the initial commits exist
    And the initial branch hierarchy exists
