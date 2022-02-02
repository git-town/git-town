Feature: cannot prepend perennial branches

  Scenario: on main branch
    Given my repo has a feature branch "feature"
    And my repo contains the commits
      | BRANCH  | LOCATION      | MESSAGE     |
      | feature | local, remote | good commit |
    And I am on the "main" branch
    When I run "git-town prepend new-branch"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      the branch "main" is not a feature branch. Only feature branches can have parent branches
      """
    And I am still on the "main" branch

  Scenario: on other perennial branch
    Given my repo has the perennial branches "qa" and "production"
    And I am on the "production" branch
    When I run "git-town prepend new-parent"
    Then it runs the commands
      | BRANCH     | COMMAND                  |
      | production | git fetch --prune --tags |
    And it prints the error:
      """
      the branch "production" is not a feature branch. Only feature branches can have parent branches
      """
    And I am still on the "production" branch
    And Git Town now has no branch hierarchy information
