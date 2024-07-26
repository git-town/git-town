Feature: parking multiple branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS |
      | feature      | feature      | main   | local     |
      | contribution | contribution |        | local     |
      | observed     | observed     | main   | local     |
      | prototype    | prototype    | main   | local     |
    And an uncommitted file
    When I run "git-town park feature contribution observed prototype"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "feature" is now parked
      """
    And branch "feature" is now parked
    And branch "contribution" is now parked
    And there are now no contribution branches
    And branch "observed" is now parked
    And there are now no observed branches
    And branch "prototype" is now parked
    And branch "prototype" is still prototype
    And the current branch is still "main"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND       |
      | main   | git add -A    |
      |        | git stash     |
      |        | git stash pop |
    And there are now no parked branches
    And the current branch is still "main"
    And the uncommitted file still exists
