Feature: in a subfolder on the main branch

  Background:
    Given the commits
      | BRANCH | LOCATION | MESSAGE       | FILE NAME        |
      | main   | local    | folder commit | new_folder/file1 |
    And the current branch is "main"
    When I run "git-town hack new" in the "new_folder" folder

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git rebase origin/main   |
      |        | git push                 |
      |        | git checkout -b new      |
    And the current branch is now "new"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | folder commit |
      | new    | local         | folder commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | new    | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND           |
      | new    | git checkout main |
      | main   | git branch -D new |
    And the current branch is now "main"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | folder commit |
    And no lineage exists now
