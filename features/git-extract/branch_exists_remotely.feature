Feature: git extract: errors when the branch exists remotely

  (see ../branch_exists_locally.feature)


  Background:
    Given I have a feature branch named "feature"
    And my coworker has a feature branch named "existing-feature"
    And I am on the "feature" branch


  Scenario: result
    Given I have an uncommitted file
    When I run `git extract existing-feature`
    Then it runs the Git commands
      | BRANCH  | COMMAND           |
      | feature | git fetch --prune |
    And I get the error "A branch named 'existing-feature' already exists"
    And I am still on the "feature" branch
    And I still have my uncommitted file

