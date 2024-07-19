Feature: Cannot create proposals for perennial branches

  Background:
    Given a Git repo clone
    And the branch
      | NAME      | TYPE      | LOCATIONS     |
      | perennial | perennial | local, origin |
    And the current branch is "perennial"
    And the origin is "git@github.com:git-town/git-town.git"
    When I run "git-town propose"

  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND                  |
      | perennial | git fetch --prune --tags |
    And it prints the error:
      """
      cannot propose perennial branches
      """
