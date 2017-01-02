Feature: Appending a branch to a perennial branch

  As a developer working on a perennial branch and coming across a number of changes I want to commit independently
  I want to be able to create a feature branch as the direct child of the perennial branch
  So that I can review and commit the changes separately without losing access to the other changes in my feature branch.


  Background:
    Given I have perennial branches named "qa" and "production"
    And the following commit exists in my repository
      | BRANCH     | LOCATION | MESSAGE           |
      | production | remote   | production_commit |
    And I am on the "production" branch
    And I have an uncommitted file
    When I run `git town-append new-child`


  Scenario: result
    Then it runs the commands
      | BRANCH     | COMMAND                              |
      | production | git fetch --prune                    |
      |            | git add -A                           |
      |            | git stash                            |
      |            | git checkout -b new-child production |
      | new-child  | git push -u origin new-child         |
      |            | git stash pop                        |
    And I end up on the "new-child" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH     | LOCATION | MESSAGE           |
      | production | remote   | production_commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH    | PARENT     |
      | new-child | production |


  Scenario: Undo
    When I run `git town-append --undo`
    Then it runs the commands
        | BRANCH     | COMMAND                    |
        | new-child  | git add -A                 |
        |            | git stash                  |
        |            | git push origin :new-child |
        |            | git checkout production    |
        | production | git branch -d new-child    |
        |            | git stash pop              |
    And I end up on the "production" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH     | LOCATION | MESSAGE           |
      | production | remote   | production_commit |

