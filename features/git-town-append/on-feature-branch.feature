Feature: Appending a branch to a feature branch

  As a developer working on a feature branch and coming across a number of changes I want to commit independently
  I want to be able to create a feature branch as the direct child of my current feature branch
  So that I can review and commit the changes separately without losing access to the other changes in my feature branch.


  Background:
    Given I have a feature branch named "existing-feature"
    And the following commits exist in my repository
      | BRANCH           | LOCATION         | MESSAGE                 | FILE NAME             | FILE CONTENT             |
      | existing-feature | local and remote | existing_feature_commit | existing_feature_file | existing feature content |
    And I am on the "existing-feature" branch
    And I have an uncommitted file


  Scenario: inserting a branch into the branch ancestry
    When I run `gt append new-child`
    Then it runs the commands
      | BRANCH           | COMMAND                                     |
      | existing-feature | git fetch --prune                           |
      |                  | git add -A                                  |
      |                  | git stash                                   |
      |                  | git checkout main                           |
      | main             | git rebase origin/main                      |
      |                  | git checkout existing-feature               |
      | existing-feature | git merge --no-edit origin/existing-feature |
      |                  | git merge --no-edit main                    |
      |                  | git branch new-child existing-feature       |
      |                  | git checkout new-child                      |
      | new-child        | git push -u origin new-child                |
      |                  | git stash pop                               |
    And I end up on the "new-child" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH           | LOCATION         | MESSAGE                 |
      | existing-feature | local and remote | existing_feature_commit |
      | new-child        | local and remote | existing_feature_commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH           | PARENT           |
      | existing-feature | main             |
      | new-child        | existing-feature |


  Scenario: Undo
    Given I run `gt append new-child`
    When I run `gt append --undo`
    Then it runs the commands
      | BRANCH           | COMMAND                       |
      | new-child        | git add -A                    |
      |                  | git stash                     |
      |                  | git push origin :new-child    |
      |                  | git checkout existing-feature |
      | existing-feature | git branch -D new-child       |
      |                  | git checkout main             |
      | main             | git checkout existing-feature |
      | existing-feature | git stash pop                 |
    And I end up on the "existing-feature" branch
    And I still have my uncommitted file
    And I am left with my original commits
    And Git Town is now aware of this branch hierarchy
      | BRANCH           | PARENT |
      | existing-feature | main   |
