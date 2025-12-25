Feature: rename a contribution branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS     |
      | contribution | contribution | main   | local, origin |
    And the commits
      | BRANCH       | LOCATION      | MESSAGE               |
      | contribution | local, origin | somebody elses commit |
    And local Git setting "color.ui" is "always"
    And the current branch is "contribution"
    When I run "git-town rename contribution new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH       | COMMAND                            |
      | contribution | git fetch --prune --tags           |
      |              | git branch --move contribution new |
      |              | git checkout new                   |
      | new          | git push -u origin new             |
      |              | git push origin :contribution      |
    And this lineage exists now
      """
      main
        new
      """
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE               |
      | new    | local, origin | somebody elses commit |
    And branch "new" still has type "contribution"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH       | COMMAND                                                   |
      | new          | git branch contribution {{ sha 'somebody elses commit' }} |
      |              | git push -u origin contribution                           |
      |              | git checkout contribution                                 |
      | contribution | git branch -D new                                         |
      |              | git push origin :new                                      |
    And the initial branches and lineage exist now
    And branch "contribution" still has type "contribution"
    And the initial commits exist now
