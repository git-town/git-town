Feature: git town-hack: requires a branch name

  As a developer forgetting to provide the name of the new branch to be created
  I should be reminded that I have to provide the branch name to this command
  So that I can use it correctly without having to look that fact up in the readme.


  Background:
    Given I have a feature branch named "existing-feature"
    And I am on the "existing-feature" branch
    And I have an uncommitted file
    When I run `git-town hack`


  Scenario: result
    Then it runs no commands
    And I get the error "no branch name provided"
    And I get the error
      """
      Usage:
        git-town hack <branch> [flags]
      """
    And I am still on the "existing-feature" branch
    And I still have my uncommitted file
