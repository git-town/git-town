Feature: git-hack: errors when the branch already exists remotely (with open changes)

  As a developer trying to create a feature branch that already exists remotely
  I should see an error telling me that a branch with that name already exists
  So that my new feature branch is unique.


  Background:
    Given my coworker has a feature branch named "existing_feature"
    And I am on the main branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git hack existing_feature` while allowing errors


  Scenario: result
    Then it runs the Git commands
      | BRANCH | COMMAND           |
      | main   | git fetch --prune |
    And I get the error "A branch named 'existing_feature' already exists"
    And I am still on the "main" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
