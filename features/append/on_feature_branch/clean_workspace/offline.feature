Feature: append in offline mode

  Background:
    Given offline mode is enabled
    And the current branch is a feature branch "existing"
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | existing | local, origin | existing commit |

  Scenario: result
    When I run "git-town append new"
    Then it runs the commands
      | BRANCH   | COMMAND                             |
      | existing | git checkout main                   |
      | main     | git rebase origin/main              |
      |          | git checkout existing               |
      | existing | git merge --no-edit origin/existing |
      |          | git merge --no-edit main            |
      |          | git checkout -b new existing        |
    And the current branch is now "new"
    And these commits exist now
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
    And the current branch is now "existing"
    And the initial commits exist
    And the initial lineage exists
