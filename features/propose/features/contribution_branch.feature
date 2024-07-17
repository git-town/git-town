Feature: Cannot create proposals for contribution branches

  Background:
    Given a Git repo clone
    And the branches
      | NAME         | TYPE         | LOCATIONS     |
      | contribution | contribution | local, origin |
    Given the current branch is "contribution"
    And the origin is "git@github.com:git-town/git-town.git"
    When I run "git-town propose"

  Scenario: result
    Then it runs the commands
      | BRANCH       | COMMAND                  |
      | contribution | git fetch --prune --tags |
    And it prints the error:
      """
      cannot propose contribution branches
      """
