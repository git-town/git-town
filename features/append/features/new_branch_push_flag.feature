Feature: auto-push the new branch to the remote

  Background:
    Given the new-branch-push-flag configuration is true
    And my repo contains the commits
      | BRANCH | LOCATION | MESSAGE     |
      | main   | remote   | main_commit |
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run "git-town append new-child"

  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND                      |
      | main      | git fetch --prune --tags     |
      |           | git add -A                   |
      |           | git stash                    |
      |           | git rebase origin/main       |
      |           | git branch new-child main    |
      |           | git checkout new-child       |
      | new-child | git push -u origin new-child |
      |           | git stash pop                |
    And I am now on the "new-child" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH    | LOCATION      | MESSAGE     |
      | main      | local, remote | main_commit |
      | new-child | local, remote | main_commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH    | PARENT |
      | new-child | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH    | COMMAND                    |
      | new-child | git add -A                 |
      |           | git stash                  |
      |           | git push origin :new-child |
      |           | git checkout main          |
      | main      | git branch -d new-child    |
      |           | git stash pop              |
    And I am now on the "main" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, remote | main_commit |
    And Git Town now has no branch hierarchy information
