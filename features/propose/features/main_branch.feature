Feature: Cannot create proposals for the main branch

  Background:
    Given a Git repo clone
    And the origin is "git@github.com:git-town/git-town.git"
    When I run "git-town propose"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      cannot propose the main branch
      """
