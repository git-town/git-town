Feature: git extract: errors when the branch exists remotely

  (see ../branch_exists_locally.feature)


  Background:
    Given I have a feature branch named "feature"
    And my coworker has a feature branch named "existing-feature"
    And I am on the "feature" branch


  Scenario: with open changes
    Given I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git extract existing-feature` it errors
    Then it runs the Git commands
      | BRANCH  | COMMAND           |
      | feature | git fetch --prune |
    And I get the error "A branch named 'existing-feature' already exists"
    And I am still on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: without open changes
    When I run `git extract existing-feature` it errors
    Then it runs the Git commands
      | BRANCH  | COMMAND           |
      | feature | git fetch --prune |
    And I get the error "A branch named 'existing-feature' already exists"
    And I am still on the "feature" branch
