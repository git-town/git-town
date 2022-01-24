Feature: cannot prepend perennial branches

  Scenario: on main branch
    Given my repo has a feature branch named "feature"
    And the following commits exist in my repo
      | BRANCH  | LOCATION      | MESSAGE     |
      | feature | local, remote | good commit |
    And I am on the "main" branch
    Given my workspace has an uncommitted file
    When I run "git-town prepend new-branch"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      the branch "main" is not a feature branch. Only feature branches can have parent branches
      """
    And I am still on the "main" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
      | remote     | main, feature |
    And my repo is left with my original commits
    And Git Town is now aware of this branch hierarchy
      | BRANCH  | PARENT |
      | feature | main   |

  Scenario: on perennial branch
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
