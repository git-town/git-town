Feature: rename the current branch without pre-push hook

  Background:
    Given the current branch is a feature branch "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | old    | local, origin | old commit  |

  Scenario: set to "false"
    Given setting "push-hook" is "false"
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
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | new    | local, origin | old commit  |

  Scenario: set to "true"
    Given setting "push-hook" is "true"
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
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | new    | local, origin | old commit  |
