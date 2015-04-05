Feature: git hack: errors when the branch exists locally

  As a developer trying to create a branch with the name of an existing branch
  I should see an error telling me that a branch with that name already exists
  So that my new feature branch is unique.


  Background:
    Given I have a feature branch named "existing_feature"
    And I am on the "main" branch


  Scenario: with open changes
    And I have an uncommitted file
    When I run `git hack existing_feature`
    Then it runs the Git commands
      | BRANCH | COMMAND           |
      | main   | git fetch --prune |
    And I get the error "A branch named 'existing_feature' already exists"
    And I am still on the "main" branch
    And I still have my uncommitted file


  Scenario: without open changes
    When I run `git hack existing_feature`
    Then it runs the Git commands
      | BRANCH | COMMAND           |
      | main   | git fetch --prune |
    And I get the error "A branch named 'existing_feature' already exists"
    And I am still on the "main" branch
