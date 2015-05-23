Feature: git extract: errors if there are not extractable commits

  As a developer trying to extract commits from a branch that has no extractable commits
  I should see an error telling me that there are no extractable commits
  So that I know when to use this command.


  Background:
    Given I have a feature branch named "feature"
    And I am on the "feature" branch


  Scenario: result
    Given I have an uncommitted file
    When I run `git extract refactor`
    Then it runs the Git commands
      | BRANCH  | COMMAND           |
      | feature | git fetch --prune |
    And I get the error "The branch 'feature' has no extractable commits."
    And I am still on the "feature" branch
    And I still have my uncommitted file

