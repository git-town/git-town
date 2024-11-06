Feature: Cannot create proposals for perennial branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | LOCATIONS     |
      | perennial | perennial | local, origin |
    And the current branch is "perennial"
    And the origin is "git@github.com:git-town/git-town.git"
    When I run "git-town propose"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                  |
      | perennial | git fetch --prune --tags |
    And Git Town prints the error:
      """
      cannot propose perennial branches
      """
