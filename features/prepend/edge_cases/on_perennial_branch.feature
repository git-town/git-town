Feature: does not prepend perennial branches

  Scenario: on main branch
    And the current branch is "main"
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
    And the current branch is a perennial branch "production"
    When I run "git-town prepend new"
    Then it runs the commands
      | BRANCH     | COMMAND                  |
      | production | git fetch --prune --tags |
    And it prints the error:
      """
      the branch "production" is not a feature branch. Only feature branches can have parent branches
      """
    And the current branch is still "production"
