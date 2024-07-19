Feature: observing multiple branches

  Background:
    Given a Git repo clone
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS     |
      | feature      | feature      | main   | local, origin |
      | contribution | contribution |        | local, origin |
      | parked       | parked       | main   | local, origin |
      | prototype    | prototype    | main   | local, origin |
    And an uncommitted file
    When I run "git-town observe feature contribution parked prototype"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "feature" is now an observed branch
      """
    And branch "feature" is now observed
    And it prints:
      """
      branch "contribution" is now an observed branch
      """
    And branch "contribution" is now observed
    And there are now no contribution branches
    And it prints:
      """
      branch "parked" is now an observed branch
      """
    And branch "parked" is now observed
    And there are now no parked branches
    And it prints:
      """
      branch "prototype" is now an observed branch
      """
    And branch "prototype" is now observed
    And there are now no prototype branches
    And the current branch is still "main"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND       |
      | main   | git add -A    |
      |        | git stash     |
      |        | git stash pop |
    And there are now no observed branches
    And the initial branches exist
    And the current branch is still "main"
    And the uncommitted file still exists
