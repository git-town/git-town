Feature: Trying to create a feature branch with a non-existing parent

  As a developer trying to branch off a non-existing branch
  I want to get a reminder about my mistake
  So that I can try again with the correct parent branch name.


  Background:
    Given I am on the "main" branch
    And I have an uncommitted file
    When I run `git hack feature zonk`


  Scenario: result
    Then I get the error
      """
      There is no branch named 'zonk'
      """
    And it runs no commands
    And I end up on the "main" branch
    And I still have my uncommitted file
    And I am left with my original commits
    And Git Town is not aware of any branch hierarchy
