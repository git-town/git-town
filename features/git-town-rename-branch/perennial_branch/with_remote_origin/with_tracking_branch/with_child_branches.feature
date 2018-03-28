Feature: git town-rename-branch: renaming a feature branch with child branches

  As a developer renaming a feature branch that has child branches
  I want that the branch hierarchy information is updated to the new branch name
  So that my workspace is in a consistent and fully functional state after the rename.


  Background:
    Given my repository has a perennial branch named "production"
    And my repository has a feature branch named "child-feature" as a child of "production"
    And the following commits exist in my repository
      | BRANCH        | LOCATION         | MESSAGE              | FILE NAME          | FILE CONTENT          |
      | child-feature | local and remote | child feature commit | child_feature_file | child feature content |
      | production    | local and remote | production commit    | production_file    | production content    |
    And I am on the "production" branch
    When I run `git-town rename-branch --force production renamed-production`


  Scenario: result
    Then it runs the commands
      | BRANCH             | COMMAND                                  |
      | production         | git fetch --prune                        |
      |                    | git branch renamed-production production |
      |                    | git checkout renamed-production          |
      | renamed-production | git push -u origin renamed-production    |
      |                    | git push origin :production              |
      |                    | git branch -D production                 |
    And I end up on the "renamed-production" branch
    And the perennial branches are now configured as "renamed-production"
    And my repository has the following commits
      | BRANCH             | LOCATION         | MESSAGE              | FILE NAME          | FILE CONTENT          |
      | child-feature      | local and remote | child feature commit | child_feature_file | child feature content |
      | renamed-production | local and remote | production commit    | production_file    | production content    |
    And Git Town is now aware of this branch hierarchy
      | BRANCH        | PARENT             |
      | child-feature | renamed-production |


  Scenario: undo
    When I run `git-town undo`
    Then it runs the commands
      | BRANCH             | COMMAND                                              |
      | renamed-production | git branch production <%= sha 'production commit' %> |
      |                    | git push -u origin production                        |
      |                    | git push origin :renamed-production                  |
      |                    | git checkout production                              |
      | production         | git branch -D renamed-production                     |
    And I end up on the "production" branch
    And the perennial branches are now configured as "production"
    And my repository has the following commits
      | BRANCH        | LOCATION         | MESSAGE              | FILE NAME          | FILE CONTENT          |
      | child-feature | local and remote | child feature commit | child_feature_file | child feature content |
      | production    | local and remote | production commit    | production_file    | production content    |
    And Git Town is now aware of this branch hierarchy
      | BRANCH        | PARENT     |
      | child-feature | production |
