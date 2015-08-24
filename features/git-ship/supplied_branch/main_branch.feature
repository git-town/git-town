Feature: git ship: errors when trying to ship the main branch

  (see ../current_branch/on_main_branch.feature)


  Background:
    Given I have a feature branch named "feature"
    And I am on the "feature" branch
    And I have an uncommitted file
    When I run `git ship main`


  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND           |
      | feature | git fetch --prune |
    And I get the error "The branch 'main' is not a feature branch. Only feature branches can be shipped."
    And I am still on the "feature" branch
    And I still have my uncommitted file
