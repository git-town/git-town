Feature: rename the current branch without pre-push hook

  Background:
    Given a Git repo clone
    And the branch
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the current branch is "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | old    | local, origin | old commit  |

  Scenario: set to "false"
    Given Git Town setting "push-hook" is "false"
    When I run "git-town rename-branch new"
    Then it runs the commands
      | BRANCH | COMMAND                            |
      | old    | git fetch --prune --tags           |
      |        | git branch new old                 |
      |        | git checkout new                   |
      | new    | git push --no-verify -u origin new |
      |        | git push origin :old               |
      |        | git branch -D old                  |
    And the current branch is now "new"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | new    | local, origin | old commit  |

  Scenario: set to "true"
    Given Git Town setting "push-hook" is "true"
    When I run "git-town rename-branch new"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | old    | git fetch --prune --tags |
      |        | git branch new old       |
      |        | git checkout new         |
      | new    | git push -u origin new   |
      |        | git push origin :old     |
      |        | git branch -D old        |
    And the current branch is now "new"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | new    | local, origin | old commit  |
