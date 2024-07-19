Feature: does not prepend perennial branches

  Background:
    Given a Git repo clone

  Scenario: on main branch
    When I run "git-town prepend new"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      the branch "main" is not a feature branch. Only feature branches can have parent branches
      """
    And the current branch is still "main"

  Scenario: on perennial branch
    Given the branch
      | NAME       | TYPE      | LOCATIONS     |
      | production | perennial | local, origin |
    And the current branch is "production"
    When I run "git-town prepend new"
    Then it runs the commands
      | BRANCH     | COMMAND                  |
      | production | git fetch --prune --tags |
    And it prints the error:
      """
      the branch "production" is not a feature branch. Only feature branches can have parent branches
      """
    And the current branch is still "production"
