Feature: make a known remote branch a contribution branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME           | TYPE    | PARENT | LOCATIONS |
      | remote-feature | feature | main   | origin    |
    And I run "git fetch"
    And an uncommitted file
    When I run "git-town contribute remote-feature"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                     |
      |        | git checkout remote-feature |
    And it prints:
      """
      branch "remote-feature" is now a contribution branch
      """
    And branch "remote-feature" is now a contribution branch
    And the current branch is now "remote-feature"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH         | COMMAND       |
      | remote-feature | git add -A    |
      |                | git stash     |
      |                | git stash pop |
    And the current branch is still "remote-feature"
    And there are now no contribution branches
    And the uncommitted file still exists
