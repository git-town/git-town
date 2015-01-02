Feature: git kill: don't delete a non-existing branch (without open changes)

  (see ./non_existing_branch_with_open_changes.feature)


  Background:
    Given I am on the "main" branch
    When I run `git kill non-existing-feature` while allowing errors

  Scenario: result
    Then it runs the Git commands
      | BRANCH | COMMAND           |
      | main   | git fetch --prune |
    And I get the error "There is no branch named 'non-existing-feature'"
    And I end up on the "main" branch
