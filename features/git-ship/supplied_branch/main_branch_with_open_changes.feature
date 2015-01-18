Feature: git ship: don't ship the main branch (with open changes)

  (see ../current_branch/on_main_branch.feature)


  Background:
    Given I have a feature branch named "feature"
    And I am on the "feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git ship main`, it errors


  Scenario: result
    Then it runs no Git commands
    And I get the error "The branch 'main' is not a feature branch. Only feature branches can be shipped."
    And I am still on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
