Feature: git town-rename-branch: renaming a feature branch with child branches


  Background:
    Given my repo has a feature branch named "parent-feature"
    And my repo has a feature branch named "child-feature" as a child of "parent-feature"
    And the following commits exist in my repo
      | BRANCH         | LOCATION      | MESSAGE               | FILE NAME           | FILE CONTENT           |
      | child-feature  | local, remote | child feature commit  | child_feature_file  | child feature content  |
      | parent-feature | local, remote | parent feature commit | parent_feature_file | parent feature content |
    And I am on the "parent-feature" branch
    When I run "git-town rename-branch parent-feature renamed-parent-feature"

  Scenario: result
    Then it runs the commands
      | BRANCH                 | COMMAND                                          |
      | parent-feature         | git fetch --prune --tags                         |
      |                        | git branch renamed-parent-feature parent-feature |
      |                        | git checkout renamed-parent-feature              |
      | renamed-parent-feature | git push -u origin renamed-parent-feature        |
      |                        | git push origin :parent-feature                  |
      |                        | git branch -D parent-feature                     |
    And I am now on the "renamed-parent-feature" branch
    And my repo now has the following commits
      | BRANCH                 | LOCATION      | MESSAGE               | FILE NAME           | FILE CONTENT           |
      | child-feature          | local, remote | child feature commit  | child_feature_file  | child feature content  |
      | renamed-parent-feature | local, remote | parent feature commit | parent_feature_file | parent feature content |
    And Git Town is now aware of this branch hierarchy
      | BRANCH                 | PARENT                 |
      | child-feature          | renamed-parent-feature |
      | renamed-parent-feature | main                   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH                 | COMMAND                                                     |
      | renamed-parent-feature | git branch parent-feature {{ sha 'parent feature commit' }} |
      |                        | git push -u origin parent-feature                           |
      |                        | git push origin :renamed-parent-feature                     |
      |                        | git checkout parent-feature                                 |
      | parent-feature         | git branch -D renamed-parent-feature                        |
    And I am now on the "parent-feature" branch
    And my repo now has the following commits
      | BRANCH         | LOCATION      | MESSAGE               | FILE NAME           | FILE CONTENT           |
      | child-feature  | local, remote | child feature commit  | child_feature_file  | child feature content  |
      | parent-feature | local, remote | parent feature commit | parent_feature_file | parent feature content |
    And Git Town is now aware of this branch hierarchy
      | BRANCH         | PARENT         |
      | child-feature  | parent-feature |
      | parent-feature | main           |
