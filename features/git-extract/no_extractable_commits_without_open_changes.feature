Feature: git-extract: errors if there are not extractable commits

  As a developer trying to extract commits from a branch that has no extractable commits
  I should see an error telling me that there are no extractable commits
  So that I know how to use this command correctly.


  Background:
    Given I have a feature branch named "feature"
    And I am on the "feature" branch
    When I run `git extract refactor` while allowing errors


  Scenario: result
    Then it runs no Git commands
    And I get the error "The branch 'feature' has no extractable commits."
    And I am still on the "feature" branch
