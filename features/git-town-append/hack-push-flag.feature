Feature: push branch to remote upon creation

  (see ../git-town-hack/hack_push_flag.feature)


  Background:
    Given the "new-branch-push-flag" configuration is set to "true"
    And the following commit exists in my repository
      | BRANCH | LOCATION | MESSAGE     |
      | main   | remote   | main_commit |
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run `git-town append new-child`


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
    And my workspace still contains my uncommitted file
    And my repository has the following commits
      | BRANCH    | LOCATION         | MESSAGE     |
      | main      | local and remote | main_commit |
      | new-child | local and remote | main_commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH    | PARENT |
      | new-child | main   |


  Scenario: Undo
    When I run `git-town undo`
    Then it runs the commands
        | BRANCH    | COMMAND                    |
        | new-child | git add -A                 |
        |           | git stash                  |
        |           | git push origin :new-child |
        |           | git checkout main          |
        | main      | git branch -d new-child    |
        |           | git stash pop              |
    And I end up on the "main" branch
    And my workspace still contains my uncommitted file
    And my repository has the following commits
      | BRANCH | LOCATION         | MESSAGE     |
      | main   | local and remote | main_commit |
