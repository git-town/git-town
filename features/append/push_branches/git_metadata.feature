Feature: disable pushing through Git metadata

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | branch-1 | feature | main   | local     |
    And Git setting "git-town.push-branches" is "false"
    And the current branch is "branch-1"
    When I run "git-town append branch-2"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                  |
      | branch-1 | git fetch --prune --tags |
      |          | git checkout -b branch-2 |
    And this lineage exists now
      """
      main
        branch-1
          branch-2
      """
    And the branches are now
      | REPOSITORY | BRANCHES                 |
      | local      | main, branch-1, branch-2 |
      | origin     | main                     |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                |
      | branch-2 | git checkout branch-1  |
      | branch-1 | git branch -D branch-2 |
    And the initial lineage exists now
    And the branches are now
      | REPOSITORY | BRANCHES       |
      | local      | main, branch-1 |
      | origin     | main           |
