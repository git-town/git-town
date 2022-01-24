Feature: git town-rename-branch: errors when the destination branch exists remotely


  (see ./destination_branch_exists_locally.feature)

  Background:
    Given my repo has a feature branch named "current-feature"
    And my coworker has a feature branch named "existing-feature"
    And the following commits exist in my repo
      | BRANCH           | LOCATION      | MESSAGE                 |
      | current-feature  | local, remote | current-feature commit  |
      | existing-feature | remote        | existing-feature commit |
    And I am on the "current-feature" branch
    And my workspace has an uncommitted file
    When I run "git-town rename-branch current-feature existing-feature"

  Scenario: result
    Then it runs the commands
      | BRANCH          | COMMAND                  |
      | current-feature | git fetch --prune --tags |
    And it prints the error:
      """
      a branch named "existing-feature" already exists
      """
    And I am still on the "current-feature" branch
    And my workspace still contains my uncommitted file
    And my repo is left with my original commits
