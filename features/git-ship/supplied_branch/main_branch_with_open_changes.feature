Feature: git ship: don't ship the main branch (with open changes)

  As a developer accidentally trying to ship the main branch
  I should be notified about my mistake
  So that I can ship the correct branch and remain productive.


  Background:
    Given I am on the "feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git ship main -m 'feature done'` while allowing errors


  Scenario: result
    Then I get the error "The branch 'main' is not a feature branch. Only feature branches can be shipped."
    And I am still on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
