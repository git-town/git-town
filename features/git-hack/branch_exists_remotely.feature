Feature: git hack: errors when the branch exists remotely

  (see ./branch_exists_locally.feature)


  Background:
    Given my coworker has a feature branch named "existing_feature"
    And I am on the "main" branch


  Scenario: with open chanhes
    Given I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git hack existing_feature` while allowing errors
    Then it runs the Git commands
      | BRANCH | COMMAND           |
      | main   | git fetch --prune |
    And I get the error "A branch named 'existing_feature' already exists"
    And I am still on the "main" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: without open changes
    When I run `git hack existing_feature` while allowing errors
    Then it runs the Git commands
      | BRANCH | COMMAND           |
      | main   | git fetch --prune |
    And I get the error "A branch named 'existing_feature' already exists"
    And I am still on the "main" branch
