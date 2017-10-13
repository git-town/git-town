Feature: git town-hack: errors when the branch exists locally

  As a developer trying to create a branch with the name of an existing branch
  I should see an error telling me that a branch with that name already exists
  So that my new feature branch is unique.


  Background:
    Given my repository has a feature branch named "existing-feature"
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run `git-town hack existing-feature`


  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND           |
      | main   | git fetch --prune |
    And it prints the error "A branch named 'existing-feature' already exists"
    And I am still on the "main" branch
    And my workspace still contains my uncommitted file
