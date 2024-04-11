Feature: auto-push new branches

  Background:
    Given Git Town setting "push-new-branches" is "true"
    And the current branch is a feature branch "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE        |
      | old    | local, origin | feature commit |

  Scenario: set to "false"
    Given Git Town setting "push-hook" is "false"
    When I run "git-town prepend new"
    Then it runs the commands
      | BRANCH | COMMAND                            |
      | old    | git fetch --prune --tags           |
      |        | git checkout main                  |
      | main   | git rebase origin/main             |
      |        | git checkout old                   |
      | old    | git merge --no-edit origin/old     |
      |        | git merge --no-edit main           |
      |        | git branch new main                |
      |        | git checkout new                   |
      | new    | git push --no-verify -u origin new |
    And the current branch is now "new"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE        |
      | old    | local, origin | feature commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | new    | main   |
      | old    | new    |

  Scenario: set to "true"
    Given Git Town setting "push-hook" is "true"
    When I run "git-town prepend new"
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
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE        |
      | old    | local, origin | feature commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | new    | main   |
      | old    | new    |
