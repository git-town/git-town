Feature: git kill: killing the current feature branch with a tracking branch (with open changes)

  As a developer working on a dead-end feature branch
  I want to be able to cleanly delete the current branch including open changes
  So that my workspace doesn't contain irrelevant branches and my productivity remains high.

  Background:
    Given I have feature branches named "feature" and "dead-feature"
    And the following commits exist in my repository
      | BRANCH       | LOCATION         | MESSAGE         | FILE NAME        |
      | dead-feature | local and remote | dead-end commit | unfortunate_file |
      | feature      | local and remote | good commit     | good_file        |
    And I am on the "dead-feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git kill`


  Scenario: result
    Then it runs the Git commands
      | BRANCH       | COMMAND                             |
      | dead-feature | git fetch --prune                   |
      |              | git add -A                          |
      |              | git commit -m 'WIP on dead-feature' |
      |              | git checkout main                   |
      | main         | git push origin :dead-feature       |
      |              | git branch -D dead-feature          |
    And I end up on the "main" branch
    And I don't have any uncommitted files
    And the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
      | remote     | main, feature |
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE     | FILE NAME |
      | feature | local and remote | good commit | good_file |


  Scenario: undoing the kill
    When I run `git kill --undo`
    Then it runs the Git commands
      | BRANCH       | COMMAND                                                  |
      | main         | git branch dead-feature <%= sha 'WIP on dead-feature' %> |
      |              | git push -u origin dead-feature                          |
      |              | git checkout dead-feature                                |
      | dead-feature | git reset <%= sha 'dead-end commit' %>                   |
      |              | git push -f origin dead-feature                          |
    And I end up on the "dead-feature" branch
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the existing branches are
      | REPOSITORY | BRANCHES                    |
      | local      | main, dead-feature, feature |
      | remote     | main, dead-feature, feature |
    And I have the following commits
      | BRANCH       | LOCATION         | MESSAGE         | FILE NAME        |
      | dead-feature | local and remote | dead-end commit | unfortunate_file |
      | feature      | local and remote | good commit     | good_file        |
