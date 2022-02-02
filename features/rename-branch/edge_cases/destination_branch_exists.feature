Feature: destination branch exists

  Scenario: destination branch exists locally
    Given my repo has the feature branches "current-feature" and "existing-feature"
    And the following commits exist in my repo
      | BRANCH           | LOCATION      | MESSAGE                 |
      | current-feature  | local, remote | current-feature commit  |
      | existing-feature | local, remote | existing-feature commit |
    And I am on the "current-feature" branch
    When I run "git-town rename-branch current-feature existing-feature"
    Then it runs the commands
      | BRANCH          | COMMAND                  |
      | current-feature | git fetch --prune --tags |
    And it prints the error:
      """
      a branch named "existing-feature" already exists
      """
    And I am still on the "current-feature" branch
    And Git Town is still aware of this branch hierarchy
      | BRANCH           | PARENT |
      | current-feature  | main   |
      | existing-feature | main   |

  Scenario: destination branch exists remotely
    Given my repo has a feature branch "current-feature"
    And my coworker has a feature branch "existing-feature"
    And the following commits exist in my repo
      | BRANCH           | LOCATION      | MESSAGE                 |
      | current-feature  | local, remote | current-feature commit  |
      | existing-feature | remote        | existing-feature commit |
    And I am on the "current-feature" branch
    When I run "git-town rename-branch current-feature existing-feature"
    Then it runs the commands
      | BRANCH          | COMMAND                  |
      | current-feature | git fetch --prune --tags |
    And it prints the error:
      """
      a branch named "existing-feature" already exists
      """
    And I am still on the "current-feature" branch
    And Git Town is still aware of this branch hierarchy
      | BRANCH          | PARENT |
      | current-feature | main   |
