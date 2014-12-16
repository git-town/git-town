Feature: git ship: don't ship unfinished work on the current branch

  As a developer in the middle of ongoing work on a feature branch
  I should be prevented from accidentally shipping an unfinished state that contains uncommitted changes
  So that my users don't experience a broken product.


  Background:
    Given I am on a feature branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git ship -m 'feature done'` while allowing errors


  Scenario: result
    Then I get the error "You cannot ship with uncommitted changes."
    And I am still on the feature branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there are no commits
