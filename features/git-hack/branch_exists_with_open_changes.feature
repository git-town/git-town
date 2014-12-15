Feature: git-hack errors when branch exists with open changes

  Background:
    Given I have a feature branch named "existing_feature"
    And I am on the main branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git hack existing_feature` while allowing errors


  Scenario: result
    Then it runs no Git commands
    And I get the error "A branch named 'existing_feature' already exists"
    And I am still on the "main" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
