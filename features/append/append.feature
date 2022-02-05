Feature: append a new feature branch to an existing feature branch

  Background:
    Given my repo has a feature branch "existing"
    And my repo contains the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | existing | local, remote | existing commit |
    And I am on the "existing" branch
    And my workspace has an uncommitted file
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
    And I am now on the "new" branch
    And my workspace still contains my uncommitted file
    And my repo now has the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | existing | local, remote | existing commit |
      | new      | local         | existing commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH   | PARENT   |
      | existing | main     |
      | new      | existing |

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
    And I am now on the "existing" branch
    And my workspace still contains my uncommitted file
    And my repo is left with my original commits
    And Git Town still has the initial branch hierarchy
