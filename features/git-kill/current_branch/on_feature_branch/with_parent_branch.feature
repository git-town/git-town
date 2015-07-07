Feature: git kill: killing the current feature branch with a parent branch

  As a developer currently working on a child feature branch that leads nowhere
  I want to be able to cleanly delete the current branch including all open changes
  So that my workspace doesn't contain irrelevant branches and my productivity remains high.


  Background:
    Given I have a feature branch named "parent-feature"
    And I have a feature branch named "child-feature" as a child of "parent-feature"
    And the following commits exist in my repository
      | BRANCH         | LOCATION         | MESSAGE               | FILE NAME           |
      | child-feature  | local and remote | child feature commit  | child_feature_file  |
      | parent-feature | local and remote | parent feature commit | parent_feature_file |
    And I am on the "child-feature" branch
    And I have an uncommitted file
    When I run `git kill`


  Scenario: result
    Then it runs the Git commands
      | BRANCH         | COMMAND                              |
      | child-feature  | git fetch --prune                    |
      |                | git add -A                           |
      |                | git commit -m 'WIP on child-feature' |
      |                | git checkout parent-feature          |
      | parent-feature | git push origin :child-feature       |
      |                | git branch -D child-feature          |
    And I end up on the "parent-feature" branch
    And I don't have any uncommitted files
    And the existing branches are
      | REPOSITORY | BRANCHES             |
      | local      | main, parent-feature |
      | remote     | main, parent-feature |
    And I have the following commits
      | BRANCH         | LOCATION         | MESSAGE               | FILE NAME           |
      | parent-feature | local and remote | parent feature commit | parent_feature_file |
    And Git Town is now aware of this branch hierarchy
      | BRANCH         | PARENT |
      | parent-feature | main   |


  Scenario: undoing the kill
    When I run `git kill --undo`
    Then it runs the Git commands
      | BRANCH         | COMMAND                                                    |
      | parent-feature | git branch child-feature <%= sha 'WIP on child-feature' %> |
      |                | git push -u origin child-feature                           |
      |                | git checkout child-feature                                 |
      | child-feature  | git reset <%= sha 'child feature commit' %>                |
      |                | git push -f origin child-feature                           |
    And I end up on the "child-feature" branch
    And I again have my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES                            |
      | local      | main, child-feature, parent-feature |
      | remote     | main, child-feature, parent-feature |
    And I have the following commits
      | BRANCH         | LOCATION         | MESSAGE               | FILE NAME           |
      | child-feature  | local and remote | child feature commit  | child_feature_file  |
      | parent-feature | local and remote | parent feature commit | parent_feature_file |
    And Git Town is now aware of this branch hierarchy
      | BRANCH         | PARENT         |
      | child-feature  | parent-feature |
      | parent-feature | main           |
