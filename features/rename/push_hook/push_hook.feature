Feature: rename the current branch without pre-push hook

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | old    | local, origin | old commit  |
    And the current branch is "old"

  Scenario: set to "false"
    Given Git setting "git-town.push-hook" is "false"
    When I run "git-town rename new"
    Then Git Town runs the commands
      | BRANCH | COMMAND                            |
      | old    | git fetch --prune --tags           |
      |        | git branch --move old new          |
      |        | git checkout new                   |
      | new    | git push --no-verify -u origin new |
      |        | git push origin :old               |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | new    | local, origin | old commit  |

  Scenario: set to "true"
    Given Git setting "git-town.push-hook" is "true"
    When I run "git-town rename new"
    Then Git Town runs the commands
      | BRANCH | COMMAND                   |
      | old    | git fetch --prune --tags  |
      |        | git branch --move old new |
      |        | git checkout new          |
      | new    | git push -u origin new    |
      |        | git push origin :old      |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | new    | local, origin | old commit  |
