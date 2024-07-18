Feature: inside an uncommitted subfolder on the current feature branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
    And the current branch is "existing"
    And a folder "uncommitted_folder"
    When I run "git-town hack new" in the "uncommitted_folder" folder

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                  |
      | existing | git fetch --prune --tags |
      |          | git checkout main        |
      | main     | git rebase origin/main   |
      |          | git checkout -b new      |
    And the current branch is now "new"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | new    | local         | main commit |
    And this lineage exists now
      | BRANCH   | PARENT |
      | existing | main   |
      | new      | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH   | COMMAND               |
      | new      | git checkout existing |
      | existing | git branch -D new     |
    And the current branch is now "existing"
    And the initial commits exist
    And the initial lineage exists
