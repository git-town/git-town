Feature: git-extract errors if there are not extractable commits

  As a developer trying to extract commits from the main branch
  I should be reminded that this command can only be run on feature branches
  So that I know how to do this correctly without having to memorize everything.


  Background:
    Given I have a feature branch named "feature"
    And I am on the "feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git extract refactor` while allowing errors


  Scenario: result
    Then it runs no Git commands
    And I get the error "The branch 'feature' has no extractable commits."
    And I am still on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
