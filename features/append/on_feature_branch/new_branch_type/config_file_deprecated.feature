Feature: append a new branch when prototype branches are configured via a deprecated setting in the config file

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | existing | local, origin | existing commit |
    And the current branch is "existing"
    And the committed configuration file:
      """
      create-prototype-branches = true
      """
    When I run "git-town append new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                  |
      | existing | git fetch --prune --tags                 |
      |          | git checkout main                        |
      | main     | git rebase origin/main --no-update-refs  |
      |          | git checkout existing                    |
      | existing | git merge --no-edit --ff main            |
      |          | git merge --no-edit --ff origin/existing |
      |          | git checkout -b new                      |
    And Git Town prints:
      """
      The Git Town configuration file contains the deprecated setting "create-prototype-branches".
      Please upgrade to the new format: create.new-branch-type = "prototype"
      """
    And the current branch is now "new"
    And branch "new" now has type "prototype"
    And the initial commits exist now
    And this lineage exists now
      | BRANCH   | PARENT   |
      | existing | main     |
      | new      | existing |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND               |
      | new      | git checkout existing |
      | existing | git branch -D new     |
    And the current branch is now "existing"
    And the initial commits exist now
    And the initial lineage exists now