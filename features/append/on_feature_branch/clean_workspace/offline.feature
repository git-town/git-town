Feature: append in offline mode

  Background:
    Given a Git repo clone
    And the branch
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | existing | local, origin | existing commit |
    And the current branch is "existing"
    And offline mode is enabled

  Scenario: result
    When I run "git-town append new"
    Then it runs the commands
      | BRANCH   | COMMAND                                  |
      | existing | git checkout main                        |
      | main     | git rebase origin/main                   |
      |          | git checkout existing                    |
      | existing | git merge --no-edit --ff origin/existing |
      |          | git merge --no-edit --ff main            |
      |          | git checkout -b new                      |
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
