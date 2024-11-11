Feature: merging a feature branch with a parked branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | parent  | parked  | main   | local     |
      | current | feature | parent | local     |
    And the current branch is "current"
    When I run "git-town merge"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                         |
      | current | git fetch --prune --tags        |
      |         | git merge --no-edit --ff parent |
      |         | git branch -D parent            |
    And the current branch is still "current"
    And this lineage exists now
      | BRANCH  | PARENT |
      | current | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                      |
      | current | git branch parent {{ sha 'initial commit' }} |
    And the current branch is still "current"
    And the initial commits exist now
    And the initial lineage exists now
