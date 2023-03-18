Feature: append a new feature branch to an existing feature branch

  Background:
    Given the current branch is a feature branch "existing"
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | existing | local, origin | existing commit |
    And an uncommitted file
    When I run "git-town append new"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                             |
      | existing | git fetch --prune --tags            |
      |          | git add -A                          |
      |          | git stash                           |
      |          | git checkout main                   |
      | main     | git rebase origin/main              |
      |          | git checkout existing               |
      | existing | git merge --no-edit origin/existing |
      |          | git merge --no-edit main            |
      |          | git branch new existing             |
      |          | git checkout new                    |
      | new      | git stash pop                       |
    And the current branch is now "new"
    And the uncommitted file still exists
    And now these commits exist
      | BRANCH   | LOCATION      | MESSAGE         |
      | existing | local, origin | existing commit |
      | new      | local         | existing commit |
    And this branch hierarchy exists now
      | BRANCH   | PARENT   |
      | existing | main     |
      | new      | existing |

  @this
  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH   | COMMAND               |
      | new      | git add -A            |
      |          | git stash             |
      |          | git checkout existing |
      | existing | git branch -D new     |
      |          | git checkout main     |
      | main     | git checkout existing |
      | existing | git stash pop         |
    And the current branch is now "existing"
    And the uncommitted file still exists
    And now the initial commits exist
    And the initial branch hierarchy exists
