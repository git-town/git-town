Feature: git ship: don't ship the main branch (without open changes)

  (see ../current_branch/on_main_branch.feature)


  Background:
    Given I am on the "feature" branch
    When I run `git ship main -m 'feature done'` while allowing errors


  Scenario: result
    Then it runs no Git commands
    And I get the error "The branch 'main' is not a feature branch. Only feature branches can be shipped."
    And I am still on the "feature" branch
