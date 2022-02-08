Feature: append in offline mode

  Background:
    Given Git Town is in offline mode
    And my repo has a feature branch "existing"
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | existing | local, remote | existing commit |
    And I am on the "existing" branch

  Scenario: result
    When I run "git-town append new"
    Then it runs the commands
      | BRANCH   | COMMAND                             |
      | existing | git checkout main                   |
      | main     | git rebase origin/main              |
      |          | git checkout existing               |
      | existing | git merge --no-edit origin/existing |
      |          | git merge --no-edit main            |
      |          | git branch new existing             |
      |          | git checkout new                    |
    And I am now on the "new" branch
    And my repo now has the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | existing | local, remote | existing commit |
      | new      | local         | existing commit |

  Scenario: undo
    Given I ran "git-town append new"
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH   | COMMAND               |
      | new      | git checkout existing |
      | existing | git branch -D new     |
      |          | git checkout main     |
      | main     | git checkout existing |
    And I am now on the "existing" branch
    And my repo is left with my initial commits
    And Git Town is now aware of the initial branch hierarchy
