Feature: git rename-branch: errors when the destination branch exists locally

  As a developer trying to rename a branch with the name of an existing branch
  I should see an error telling me that a branch with that name already exists
  So that my new feature branch is unique.


  Background:
    Given I have a feature branches named "current-feature" and "existing-feature"
    And I am on the "current-feature" branch


  Scenario: with open changes
    And I have an uncommitted file
    When I run `git rename-branch current-feature existing-feature`
    Then it runs the Git commands
      | BRANCH          | COMMAND           |
      | current-feature | git fetch --prune |
    And I get the error "A branch named 'existing-feature' already exists"
    And I am still on the "current-feature" branch
    And I still have my uncommitted file


  Scenario: without open changes
    When I run `git rename-branch current-feature existing-feature`
    Then it runs the Git commands
      | BRANCH          | COMMAND           |
      | current-feature | git fetch --prune |
    And I get the error "A branch named 'existing-feature' already exists"
    And I am still on the "current-feature" branch
