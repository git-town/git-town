Feature: inside an uncommitted subfolder on the current feature branch

  Background:
    Given a Git repo clone
    And the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
    And an uncommitted file with name "uncommitted_folder/uncommitted" and content "uncommitted"
    When I run "git-town hack new" in the "uncommitted_folder" folder

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND             |
      | main   | git add -A          |
      |        | git stash           |
      |        | git checkout -b new |
      | new    | git stash pop       |
    And the current branch is now "new"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | new    | local         | main commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | new    | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND           |
      | new    | git add -A        |
      |        | git stash         |
      |        | git checkout main |
      | main   | git branch -D new |
      |        | git stash pop     |
    And the current branch is now "main"
    And the initial commits exist
    And the initial lineage exists
