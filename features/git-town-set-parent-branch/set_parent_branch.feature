Feature: update the parent of a nested feature branch

  As a user with a nested feature branch whose parent was shipped from another machine
  I want to be able to update the parent branch for my nested feature branch
  So that I can use it without messing with the git configuration directly


  Background:
    Given my repository has a feature branch named "parent-feature"
    And my repository has a feature branch named "child-feature" as a child of "parent-feature"
    And I am on the "child-feature" branch


  Scenario: selecting the default branch (current parent)
    When I run "git-town set-parent-branch" and answer the prompts:
      | PROMPT                                              | ANSWER  |
      | Please specify the parent branch of 'child-feature' | [ENTER] |
    Then Git Town is now aware of this branch hierarchy
      | BRANCH         | PARENT         |
      | child-feature  | parent-feature |
      | parent-feature | main           |


  Scenario: selecting another branch
    When I run "git-town set-parent-branch" and answer the prompts:
      | PROMPT                                              | ANSWER      |
      | Please specify the parent branch of 'child-feature' | [UP][ENTER] |
    Then Git Town is now aware of this branch hierarchy
      | BRANCH         | PARENT |
      | child-feature  | main   |
      | parent-feature | main   |


  Scenario: choosing "<none> (make a perennial branch)"
    When I run "git-town set-parent-branch" and answer the prompts:
      | PROMPT                                              | ANSWER          |
      | Please specify the parent branch of 'child-feature' | [UP][UP][ENTER] |
    Then the perennial branches are now configured as "child-feature"
    And Git Town is now aware of this branch hierarchy
      | BRANCH         | PARENT |
      | parent-feature | main   |
