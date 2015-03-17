Feature: git kill: killing the current feature branch without a tracking branch (with open changes)

  As a developer working on a local dead-end feature branch
  I want to be able to remove the current branch including open changes
  So that my workspace doesn't contain irrelevant branches and my productivity remains high.


  Background:
    Given I have a feature branch named "feature"
    And I have a local feature branch named "dead-feature"
    And the following commits exist in my repository
      | BRANCH       | LOCATION         | MESSAGE         | FILE NAME        |
      | feature      | local and remote | good commit     | good_file        |
      | dead-feature | local            | dead-end commit | unfortunate_file |
    And I am on the "dead-feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git kill`


  Scenario: result
    Then it runs the Git commands
      | BRANCH       | COMMAND                             |
      | dead-feature | git fetch --prune                   |
      | dead-feature | git add -A                          |
      | dead-feature | git commit -m 'WIP on dead-feature' |
      | dead-feature | git checkout main                   |
      | main         | git branch -D dead-feature          |
    And I end up on the "main" branch
    And the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
      | remote     | main, feature |
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE     | FILE NAME |
      | feature | local and remote | good commit | good_file |


  Scenario: Undoing a kill of a local feature branch
    When I run `git kill --undo`
    Then it runs the Git commands
      | BRANCH       | COMMAND                                                  |
      | main         | git branch dead-feature <%= sha 'WIP on dead-feature' %> |
      | main         | git checkout dead-feature                                |
      | dead-feature | git reset <%= sha 'dead-end commit' %>                   |
    And I end up on the "dead-feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the existing branches are
      | REPOSITORY | BRANCHES                    |
      | local      | main, dead-feature, feature |
      | remote     | main, feature               |
    And I have the following commits
      | BRANCH       | LOCATION         | MESSAGE         | FILE NAME        |
      | feature      | local and remote | good commit     | good_file        |
      | dead-feature | local            | dead-end commit | unfortunate_file |
