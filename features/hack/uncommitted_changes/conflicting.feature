Feature: unmergeable conflict between uncommitted changes and the main branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE            | FILE NAME        | FILE CONTENT |
      | main   | local, origin | conflicting commit | conflicting_file | main content |
    And the current branch is "existing"
    And an uncommitted file "conflicting_file" with content "conflicting content"
    When I run "git-town hack new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                     |
      | existing | git add -A                  |
      |          | git stash -m "Git Town WIP" |
      |          | git checkout -b new main    |
      | new      | git stash pop               |
      |          | git stash drop              |
      |          | git restore --staged .      |
    And Git Town prints:
      """
      conflicts between your uncommmitted changes and the main branch
      """
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
      | BRANCH   | COMMAND                     |
      | new      | git add -A                  |
      |          | git stash -m "Git Town WIP" |
      |          | git checkout existing       |
      | existing | git branch -D new           |
      |          | git stash pop               |
      |          | git stash drop              |
      |          | git restore --staged .      |
    And file "conflicting_file" still has content:
      """
      <<<<<<< Updated upstream
      main content
      =======
      conflicting content
      >>>>>>> Stashed changes
      """
