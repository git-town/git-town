Feature: git extract: errors when the branch exists remotely (without open changes)

  (see ../branch_exists_locally/with_open_changes.feature)


  Background:
    Given I have a feature branch named "feature"
    And my coworker has a feature branch named "existing-feature"
    And I am on the "feature" branch
    When I run `git extract existing-feature` while allowing errors


  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND           |
      | feature | git fetch --prune |
    And I get the error "A branch named 'existing-feature' already exists"
    And I am still on the "feature" branch
