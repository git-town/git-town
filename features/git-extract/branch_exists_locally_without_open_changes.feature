Feature: git-extract errors when the branch exists locally (without open changes)

  (see ./branch_exists_locally_with_open_changes.feature)


  Background:
    Given I have feature branches named "feature" and "existing-feature"
    And I am on the "feature" branch
    When I run `git extract existing-feature` while allowing errors


  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND           |
      | feature | git fetch --prune |
    And I get the error "A branch named 'existing-feature' already exists"
    And I am still on the "feature" branch
