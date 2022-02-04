Feature: does not prepend perennial branches

  Scenario: on main branch
    And I am on the "main" branch
    When I run "git-town prepend feature"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      the branch "main" is not a feature branch. Only feature branches can have parent branches
      """
    And I am still on the "main" branch

  Scenario: on perennial branch
    Given my repo has a perennial branch "production"
    And I am on the "production" branch
    When I run "git-town prepend feature"
    Then it runs the commands
      | BRANCH     | COMMAND                  |
      | production | git fetch --prune --tags |
    And it prints the error:
      """
      the branch "production" is not a feature branch. Only feature branches can have parent branches
      """
    And I am still on the "production" branch
