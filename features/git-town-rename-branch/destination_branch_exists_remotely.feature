Feature: git town-rename-branch: errors when the destination branch exists remotely

  (see ./destination_branch_exists_locally.feature)


  Background:
    Given I have a feature branch named "current-feature"
    And my coworker has a feature branch named "existing-feature"
    And the following commits exist in my repository
      | BRANCH           | LOCATION         | MESSAGE                 |
      | current-feature  | local and remote | current-feature commit  |
      | existing-feature | remote           | existing-feature commit |
    And I am on the "current-feature" branch
    And I have an uncommitted file
    When I run `git town-rename-branch current-feature existing-feature`


  Scenario: result
    Then it runs the commands
      | BRANCH          | COMMAND           |
      | current-feature | git fetch --prune |
    And I get the error "A branch named 'existing-feature' already exists"
    And I am still on the "current-feature" branch
    And I still have my uncommitted file
    And I am left with my original commits
