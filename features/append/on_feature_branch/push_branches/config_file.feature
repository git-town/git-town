Feature: disable pushing through the config file

  Background:
    Given a Git repo with origin
    And Git Town is not configured
    And the committed configuration file:
      """
      [branches]
      main = "main"
      
      [sync]
      push-branches = false
      """
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | branch-1 | feature | main   | local     |
    And the current branch is "branch-1"
    When I run "git-town append branch-2"

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
