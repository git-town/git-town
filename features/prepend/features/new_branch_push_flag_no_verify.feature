Feature: auto-push new branches

  Background:
    Given setting "new-branch-push-flag" is "true"
    And the current branch is a feature branch "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE        |
      | old    | local, origin | feature commit |

  Scenario: set to "false"
    Given setting "push-verify" is "false"
    When I run "git-town prepend new"
    Then it runs the commands
      | BRANCH | COMMAND                            |
      | old    | git fetch --prune --tags           |
      |        | git checkout main                  |
      | main   | git rebase origin/main             |
      |        | git branch new main                |
      |        | git checkout new                   |
      | new    | git push --no-verify -u origin new |
    And the current branch is now "new"
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE        |
      | old    | local, origin | feature commit |
    And this branch hierarchy exists now
      | BRANCH | PARENT |
      | new    | main   |
      | old    | new    |

  Scenario: set to "true"
    Given setting "push-verify" is "true"
    When I run "git-town prepend new"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | old    | git fetch --prune --tags |
      |        | git checkout main        |
      | main   | git rebase origin/main   |
      |        | git branch new main      |
      |        | git checkout new         |
      | new    | git push -u origin new   |
    And the current branch is now "new"
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE        |
      | old    | local, origin | feature commit |
    And this branch hierarchy exists now
      | BRANCH | PARENT |
      | new    | main   |
      | old    | new    |
