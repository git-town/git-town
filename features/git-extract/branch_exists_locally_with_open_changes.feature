Feature: git extract: don't extract into an already existing branch (with open changes)

  As a developer trying to extract commits into a branch with the name of an existing branch
  I should get an error that this branch already exists
  So that all my feature branches are unique.


  Background:
    Given I have feature branches named "feature" and "existing-feature"
    And I am on the "feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git extract existing-feature` while allowing errors


  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND           |
      | feature | git fetch --prune |
    And I get the error "A branch named 'existing-feature' already exists"
    And I am still on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
