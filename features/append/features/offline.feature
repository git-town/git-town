Feature: append in offline mode

  Background:
    Given Git Town is in offline mode
    And my repo has a feature branch "existing"
    And my repo contains the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | existing | local, origin | existing commit |
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
    And now these commits exist
      | BRANCH   | LOCATION      | MESSAGE         |
      | existing | local, origin | existing commit |
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
    And now the initial commits exist
    And Git Town is now aware of the initial branch hierarchy
