Feature: make the current observed branch a contribution branch

  Background:
    Given a Git repo with origin
    And the branch
      | NAME   | TYPE     | LOCATIONS |
      | branch | observed | local     |
    And the current branch is "branch"
    When I run "git-town contribute"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "branch" is now a contribution branch
      """
    And branch "branch" is now a contribution branch
    And the current branch is still "branch"
    And there are now no observed branches

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND       |
      | branch | git add -A    |
      |        | git stash     |
      |        | git stash pop |
    And the current branch is still "branch"
    And branch "branch" is now observed
    And there are now no contribution branches
