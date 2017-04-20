Feature: Appending a branch to a feature branch

  As a developer working on a feature branch and coming across a number of changes I want to commit independently
  I want to be able to create a feature branch as the direct child of my current feature branch
  So that I can review and commit the changes separately without losing access to the other changes in my feature branch.


  Background:
    Given the following commit exists in my repository
      | BRANCH | LOCATION | MESSAGE     |
      | main   | remote   | main_commit |
    And I am on the "main" branch
    And I have an uncommitted file
    When I run `gt append new-child`


  Scenario: inserting a branch into the branch ancestry
    Then it runs the commands
      | BRANCH    | COMMAND                      |
      | main      | git fetch --prune            |
      |           | git add -A                   |
      |           | git stash                    |
      |           | git rebase origin/main       |
      |           | git branch new-child main    |
      |           | git checkout new-child       |
      | new-child | git push -u origin new-child |
      |           | git stash pop                |
    And I end up on the "new-child" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH    | LOCATION         | MESSAGE     |
      | main      | local and remote | main_commit |
      | new-child | local and remote | main_commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH    | PARENT |
      | new-child | main   |


  Scenario: Undo
    When I run `gt append --undo`
    Then it runs the commands
        | BRANCH    | COMMAND                    |
        | new-child | git add -A                 |
        |           | git stash                  |
        |           | git push origin :new-child |
        |           | git checkout main          |
        | main      | git branch -d new-child    |
        |           | git stash pop              |
    And I end up on the "main" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH | LOCATION         | MESSAGE     |
      | main   | local and remote | main_commit |
