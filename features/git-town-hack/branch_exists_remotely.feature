Feature: git town-hack: errors when the branch exists remotely

  (see ./branch_exists_locally.feature)


  Background:
    Given my coworker has a feature branch named "existing-feature"
    And I am on the "main" branch
    And I have an uncommitted file
    When I run `git town-hack existing-feature`


  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND           |
      | main   | git fetch --prune |
    And I get the error "A branch named 'existing-feature' already exists"
    And I am still on the "main" branch
    And I still have my uncommitted file
