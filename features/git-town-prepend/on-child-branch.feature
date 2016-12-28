Feature: Prepending a branch to a nested feature branch

  As a developer working on a nested feature branch and coming across a number of changes I want to commit independently
  I want to be able to insert a feature branch as the direct parent of my current feature branch
  So that I can review and commit the changes separately without losing access to them in my current feature branch.

  - running "git town-prepend" on a feature branch creates a new feature branch as a parent of the current branch
  - the user ends up on the new feature branch
  - all open changes are carried over


  Background:
    Given I have a feature branch named "existing-parent"
    And I have a feature branch named "existing-feature" as a child of "existing-parent"
    And the following commits exist in my repository
      | BRANCH           | LOCATION         | MESSAGE                 | FILE NAME             | FILE CONTENT             |
      | existing-parent  | local and remote | existing_parent_commit  | existing_parent_file  | existing parent content  |
      | existing-feature | local and remote | existing_feature_commit | existing_feature_file | existing feature content |
    And I am on the "existing-feature" branch
    And I have an uncommitted file


  Scenario: inserting a branch into the branch ancestry
    When I run `git town-prepend new-parent`
    Then it runs the commands
      | BRANCH           | COMMAND                                    |
      | existing-feature | git fetch --prune                          |
      |                  | git add -A                                 |
      |                  | git stash                                  |
      |                  | git checkout main                          |
      | main             | git rebase origin/main                     |
      |                  | git checkout existing-parent               |
      | existing-parent  | git merge --no-edit origin/existing-parent |
      |                  | git merge --no-edit main                   |
      |                  | git checkout -b new-parent existing-parent |
      | new-parent       | git push -u origin new-parent              |
      |                  | git stash pop                              |
    And I end up on the "new-parent" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH           | LOCATION         | MESSAGE                 |
      | existing-feature | local and remote | existing_feature_commit |
      | existing-parent  | local and remote | existing_parent_commit  |
      | new-parent       | local and remote | existing_parent_commit  |
    And Git Town is now aware of this branch hierarchy
      | BRANCH           | PARENT          |
      | existing-feature | new-parent      |
      | existing-parent  | main            |
      | new-parent       | existing-parent |


  Scenario: Undo
    Given I run `git town-prepend new-parent`
    When I run `git town-prepend --undo`
    Then it runs the commands
      | BRANCH          | COMMAND                      |
      | child-feature   | git add -A                   |
      |                 | git stash                    |
      |                 | git push origin :new-parent  |
      |                 | git checkout existing-parent |
      | existing-parent | git branch -d new-parent     |
      |                 | git stash pop                |
    And I end up on the "existing-parent" branch
    And I still have my uncommitted file
    And I am left with my original commits
    And Git Town is now aware of this branch hierarchy
      | BRANCH           | PARENT          |
      | existing-feature | existing-parent |
      | existing-parent  | main            |
