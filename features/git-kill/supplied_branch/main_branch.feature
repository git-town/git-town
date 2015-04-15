Feature: git kill: errors when trying to kill the main branch

  (see ../current_branch/on_main_branch.feature)


  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION         | MESSAGE     | FILE NAME |
      | feature | local and remote | good commit | good_file |
      | main    | local and remote | main commit | main_file |
    And I am on the "feature" branch


  Scenario: with open changes
    Given I have an uncommitted file
    When I run `git kill main`
    Then it runs no Git commands
    And I get the error "You can only kill feature branches"
    And I am still on the "feature" branch
    And I still have my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
      | remote     | main, feature |
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE     | FILE NAME |
      | main    | local and remote | main commit | main_file |
      | feature | local and remote | good commit | good_file |


  Scenario: without open changes
    When I run `git kill main`
    Then it runs no Git commands
    And I get the error "You can only kill feature branches"
    And I am still on the "feature" branch
    And the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
      | remote     | main, feature |
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE     | FILE NAME |
      | main    | local and remote | main commit | main_file |
      | feature | local and remote | good commit | good_file |
