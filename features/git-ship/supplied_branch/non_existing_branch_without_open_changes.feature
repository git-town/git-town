Feature: git ship: don't ship non-existing branches (without open changes)

  (see ./non_existing_branch_with_open_changes.feature)


  Background:
    Given I am on the "main" branch
    When I run `git ship non-existing-branch` it errors


  Scenario: result
    Then it runs no Git commands
    And I get the error "There is no branch named 'non-existing-branch'"
    And I end up on the "main" branch
