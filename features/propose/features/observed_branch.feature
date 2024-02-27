Feature: Cannot create proposals for observed branches

  Background:
    Given the current branch is an observed branch "observed"
    And the origin is "git@github.com:git-town/git-town.git"
    When I run "git-town propose"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                  |
      | observed | git fetch --prune --tags |
    And it prints the error:
      """
      cannot propose observed branches
      """
