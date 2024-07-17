Feature: auto-push new branches

  Background:
    Given a Git repo clone
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    Given Git Town setting "push-new-branches" is "true"
    And Git Town setting "push-hook" is "false"
    And the current branch is "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE        |
      | old    | local, origin | feature commit |
    When I run "git-town prepend new"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                             |
      | old    | git fetch --prune --tags            |
      |        | git checkout main                   |
      | main   | git rebase origin/main              |
      |        | git checkout old                    |
      | old    | git merge --no-edit --ff origin/old |
      |        | git merge --no-edit --ff main       |
      |        | git checkout -b new main            |
      | new    | git push --no-verify -u origin new  |
    And the current branch is now "new"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE        |
      | old    | local, origin | feature commit |
    And this lineage exists now
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
    And the current branch is now "old"
    And the initial commits exist
    And the initial lineage exists
