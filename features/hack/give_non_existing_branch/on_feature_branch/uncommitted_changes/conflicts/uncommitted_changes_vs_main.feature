Feature: conflicts between uncommitted changes and the main branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the current branch is "existing"
    And the commits
      | BRANCH | LOCATION      | MESSAGE            | FILE NAME        | FILE CONTENT |
      | main   | local, origin | conflicting commit | conflicting_file | main content |
    And an uncommitted file with name "conflicting_file" and content "conflicting content"
    When I run "git-town hack new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                     |
      | existing | git add -A                  |
      |          | git stash -m "Git Town WIP" |
      |          | git checkout -b new main    |
      | new      | git stash pop               |
      |          | git stash drop              |
    And file "conflicting_file" now has content:
      """
      <<<<<<< Updated upstream
      main content
      =======
      conflicting content
      >>>>>>> Stashed changes
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND               |
      | new      | git add -A            |
      |          | git stash             |
      |          | git checkout existing |
      | existing | git branch -D new     |
      |          | git stash pop         |
      |          | git stash drop        |
    And the current branch is now "existing"
    And file "conflicting_file" still has content:
      """
      <<<<<<< Updated upstream
      main content
      =======
      conflicting content
      >>>>>>> Stashed changes
      """

  Scenario: continue with unresolved conflict
    When I run "git-town continue"
    Then Git Town prints the error:
      """
      you must resolve the conflicts before continuing
      """

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue" and close the editor
    Then Git Town runs no commands
