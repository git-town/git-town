Feature: update the parent of a nested feature branch

  As a user with a nested feature branch shipped whose parent was shipped from another machine
  I want to be able to update the parent branch for my nested feature branch
  So that I can use it without messing with the git configuration directly


  Background:
    Given I have a feature branch named "parent-feature"
    And I have a feature branch named "child-feature" as a child of "parent-feature"


  Scenario: updating the parent branch
    When I run `git town parent-branch child-feature main`
    Then Git Town is now aware of this branch hierarchy
      | BRANCH         | PARENT |
      | child-feature  | main   |
      | parent-feature | main   |


  Scenario: invalid child branch name
    When I run `git town parent-branch non-existing parent-feature`
    Then I get the error
      """
      error: no branch named 'non-existing'
      """


  Scenario: invalid parent branch name
    When I run `git town parent-branch child-feature non-existing`
    Then I get the error
      """
      error: no branch named 'non-existing'
      """
