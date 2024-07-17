Feature: inside a committed subfolder that exists only on the current feature branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    Given the current branch is "existing"
    And the commits
      | BRANCH   | LOCATION      | MESSAGE       | FILE NAME              |
      | existing | local, origin | folder commit | committed_folder/file1 |
    And an uncommitted file with name "committed_folder/uncommitted" and content "uncommitted"
    When I run "git-town hack new" in the "committed_folder" folder

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                  |
      | existing | git add -A               |
      |          | git stash                |
      |          | git checkout -b new main |
      | new      | git stash pop            |
    And the current branch is now "new"
    And the initial commits exist
    And this lineage exists now
      | BRANCH   | PARENT |
      | existing | main   |
      | new      | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH   | COMMAND               |
      | new      | git add -A            |
      |          | git stash             |
      |          | git checkout existing |
      | existing | git branch -D new     |
      |          | git stash pop         |
    And the current branch is now "existing"
    And the initial commits exist
    And the initial lineage exists
