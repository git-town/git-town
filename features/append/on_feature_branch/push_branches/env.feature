Feature: avoid pushing through the environment

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | branch-1 | feature | main   | local     |
    And the current branch is "branch-1"
    When I run "git-town append branch-2" with these environment variables
      | GIT_TOWN_PUSH_BRANCHES | false |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                  |
      | branch-1 | git fetch --prune --tags |
      |          | git checkout -b branch-2 |
    And the branches are now
      | REPOSITORY | BRANCHES                 |
      | local      | main, branch-1, branch-2 |
      | origin     | main                     |
    And this lineage exists now
      """
      main
        branch-1
          branch-2
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                |
      | branch-2 | git checkout branch-1  |
      | branch-1 | git branch -D branch-2 |
    And the branches are now
      | REPOSITORY | BRANCHES       |
      | local      | main, branch-1 |
      | origin     | main           |
    And the initial lineage exists now
