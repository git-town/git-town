Feature: cannot hack when branch exists without open changes


  Background:
    Given I have a feature branch named "feature"
    And I am on the main branch
    When I run `git hack feature` while allowing errors


  Scenario: result
    Then I get the error "A branch named 'feature' already exists"
    And I am still on the "main" branch
