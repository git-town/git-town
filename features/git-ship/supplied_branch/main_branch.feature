Feature: git ship: errors when trying to ship the main branch

  (see ../current_branch/on_main_branch.feature)


  Background:
    Given I have a feature branch named "feature"
    And I am on the "feature" branch


  Scenario: with open changes
    Given I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git ship main` while allowing errors
    Then it runs no Git commands
    And I get the error "The branch 'main' is not a feature branch. Only feature branches can be shipped."
    And I am still on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: without open changes
    When I run `git ship main` while allowing errors
    Then it runs no Git commands
    And I get the error "The branch 'main' is not a feature branch. Only feature branches can be shipped."
    And I am still on the "feature" branch
