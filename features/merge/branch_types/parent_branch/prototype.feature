Feature: merging a feature branch with a prototype branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE      | PARENT | LOCATIONS |
      | parent  | prototype | main   | local     |
      | current | feature   | parent | local     |
    And the current branch is "current"
    When I run "git-town merge"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | current | git fetch --prune --tags |
      |         | git checkout parent      |
      | parent  | git branch -D current    |
    And this lineage exists now
      """
      main
        parent
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | parent | git branch current {{ sha 'initial commit' }} |
      |        | git checkout current                          |
    And the initial lineage exists now
    And the initial commits exist now
