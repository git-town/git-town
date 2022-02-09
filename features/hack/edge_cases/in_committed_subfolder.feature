Feature: inside a committed subfolder that exists only on the current feature branch

  Background:
    Given the current branch is a feature branch "existing"
    And the commits
      | BRANCH   | LOCATION      | MESSAGE       | FILE NAME        |
      | existing | local, origin | folder commit | new_folder/file1 |
    When I run "git-town hack new" in the "new_folder" folder

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                  |
      | existing | git fetch --prune --tags |
      |          | git checkout main        |
      | main     | git rebase origin/main   |
      |          | git branch new main      |
      |          | git checkout new         |
    And the current branch is now "new"
    And now the initial commits exist
    And this branch hierarchy exists now
      | BRANCH   | PARENT |
      | existing | main   |
      | new      | main   |

  Scenario: undo
    When I run "git town undo"
    Then it runs the commands
      | BRANCH | COMMAND               |
      | new    | git checkout main     |
      | main   | git branch -d new     |
      |        | git checkout existing |
    And the current branch is now "existing"
    And now the initial commits exist
    And the initial branch hierarchy exists
