Feature: Appending a branch to a parent branch

  As a developer working on a feature branch with existing children and coming across a number of changes I want to commit independently
  I want to be able to create a feature branch as the direct child of my current feature branch
  So that I can review and commit the changes separately without losing access to the other changes in my feature branch.

  - it creates a new child branch
  - existing child branches are not modified
  - run "git town-prepend" from a child to insert the branch in front of it


  Background:
    Given I have a feature branch named "existing-parent"
    And I have a feature branch named "existing-child" as a child of "existing-parent"
    And the following commits exist in my repository
      | BRANCH          | LOCATION         | MESSAGE                | FILE NAME            | FILE CONTENT            |
      | existing-parent | local and remote | existing_parent_commit | existing_parent_file | existing parent content |
    And I am on the "existing-parent" branch
    And I have an uncommitted file
    When I run `git town-append new-child`


  Scenario: inserting a branch into the branch ancestry
    Then it runs the commands
      | BRANCH          | COMMAND                                   |
      | existing-parent | git fetch --prune                         |
      |                 | git add -A                                |
      |                 | git stash                                 |
      |                 | git checkout main                         |
      | main            | git rebase origin/main                    |
      |                 | git checkout -b new-child existing-parent |
      | new-child       | git push -u origin new-child              |
      |                 | git stash pop                             |
    And I end up on the "new-child" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH          | LOCATION         | MESSAGE                |
      | existing-parent | local and remote | existing_parent_commit |
      | new-child       | local and remote | existing_parent_commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH          | PARENT          |
      | existing-child  | existing-parent |
      | existing-parent | main            |
      | new-child       | existing-parent |


  Scenario: Undo
    When I run `git town-append --undo`
    Then it runs the commands
        | BRANCH          | COMMAND                      |
        | new-child       | git add -A                   |
        |                 | git stash                    |
        |                 | git push origin :new-child   |
        |                 | git checkout main            |
        | main            | git branch -D new-child      |
        |                 | git checkout existing-parent |
        | existing-parent | git stash pop                |
    And I end up on the "existing-parent" branch
    And I still have my uncommitted file
    And I am left with my original commits
    And Git Town is now aware of this branch hierarchy
      | BRANCH          | PARENT          |
      | existing-child  | existing-parent |
      | existing-parent | main            |
