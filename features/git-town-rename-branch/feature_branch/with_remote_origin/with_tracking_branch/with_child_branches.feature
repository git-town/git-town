Feature: git town-rename-branch: renaming a feature branch with child branches

  As a developer renaming a feature branch that has child branches
  I want that the branch hierarchy information is updated to the new branch name
  So that my workspace is in a consistent and fully functional state after the rename.


  Background:
    Given I have a feature branch named "parent-feature"
    And I have a feature branch named "child-feature" as a child of "parent-feature"
    And the following commits exist in my repository
      | BRANCH         | LOCATION         | MESSAGE               | FILE NAME           | FILE CONTENT           |
      | child-feature  | local and remote | child feature commit  | child_feature_file  | child feature content  |
      | parent-feature | local and remote | parent feature commit | parent_feature_file | parent feature content |
    And I am on the "parent-feature" branch
    When I run `gt rename-branch parent-feature renamed-parent-feature`


  Scenario: result
    Then it runs the commands
      | BRANCH                 | COMMAND                                          |
      | parent-feature         | git fetch --prune                                |
      |                        | git branch renamed-parent-feature parent-feature |
      |                        | git checkout renamed-parent-feature              |
      | renamed-parent-feature | git push -u origin renamed-parent-feature        |
      |                        | git push origin :parent-feature                  |
      |                        | git branch -D parent-feature                     |
    And I end up on the "renamed-parent-feature" branch
    And I have the following commits
      | BRANCH                 | LOCATION         | MESSAGE               | FILE NAME           | FILE CONTENT           |
      | child-feature          | local and remote | child feature commit  | child_feature_file  | child feature content  |
      | renamed-parent-feature | local and remote | parent feature commit | parent_feature_file | parent feature content |
    And Git Town is now aware of this branch hierarchy
      | BRANCH                 | PARENT                 |
      | child-feature          | renamed-parent-feature |
      | renamed-parent-feature | main                   |


  Scenario: undo
    When I run `gt rename-branch --undo`
    Then it runs the commands
      | BRANCH                 | COMMAND                                                      |
      | renamed-parent-feature | git branch parent-feature <%= sha 'parent feature commit' %> |
      |                        | git push -u origin parent-feature                            |
      |                        | git push origin :renamed-parent-feature                      |
      |                        | git checkout parent-feature                                  |
      | parent-feature         | git branch -D renamed-parent-feature                         |
    And I end up on the "parent-feature" branch
    And I have the following commits
      | BRANCH         | LOCATION         | MESSAGE               | FILE NAME           | FILE CONTENT           |
      | child-feature  | local and remote | child feature commit  | child_feature_file  | child feature content  |
      | parent-feature | local and remote | parent feature commit | parent_feature_file | parent feature content |
    And Git Town is now aware of this branch hierarchy
      | BRANCH         | PARENT         |
      | child-feature  | parent-feature |
      | parent-feature | main           |
