Feature: Appending a branch to a perennial branch

  As a developer working on a perennial branch and coming across a number of changes I want to commit independently
  I want to be able to create a feature branch as the direct child of the perennial branch
  So that I can review and commit the changes separately without losing access to the other changes in my feature branch.


  Background:
    Given my repository has the perennial branches "qa" and "production"
    And the following commit exists in my repository
      | BRANCH     | LOCATION | MESSAGE           |
      | production | remote   | production_commit |
    And I am on the "production" branch
    And my workspace has an uncommitted file
    When I run `git-town append new-child`


  Scenario: result
    Then it runs the commands
      | BRANCH     | COMMAND                         |
      | production | git fetch --prune               |
      |            | git add -A                      |
      |            | git stash                       |
      |            | git rebase origin/production    |
      |            | git branch new-child production |
      |            | git checkout new-child          |
      | new-child  | git stash pop                   |
    And I end up on the "new-child" branch
    And my workspace still contains my uncommitted file
    And my repository has the following commits
      | BRANCH     | LOCATION         | MESSAGE           |
      | new-child  | local            | production_commit |
      | production | local and remote | production_commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH    | PARENT     |
      | new-child | production |


  Scenario: Undo
    When I run `git-town undo`
    Then it runs the commands
        | BRANCH     | COMMAND                 |
        | new-child  | git add -A              |
        |            | git stash               |
        |            | git checkout production |
        | production | git branch -D new-child |
        |            | git stash pop           |
    And I end up on the "production" branch
    And my workspace still contains my uncommitted file
    And my repository has the following commits
      | BRANCH     | LOCATION         | MESSAGE           |
      | production | local and remote | production_commit |
