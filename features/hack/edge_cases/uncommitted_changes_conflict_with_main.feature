Feature: conflicts between uncommitted changes and the main branch

  Background:
    Given the current branch is a feature branch "existing"
    And the commits
      | BRANCH | LOCATION      | MESSAGE            | FILE NAME        | FILE CONTENT |
      | main   | local, origin | conflicting commit | conflicting_file | main content |
    And an uncommitted file with name "conflicting_file" and content "conflicting content"
    When I run "git-town hack new"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                  |
      | existing | git fetch --prune --tags |
      |          | git add -A               |
      |          | git stash                |
      |          | git checkout main        |
      | main     | git rebase origin/main   |
      |          | git branch new main      |
      |          | git checkout new         |
      | new      | git stash pop            |
    And it prints the error:
      """
      conflicts between your uncommmitted changes and the main branch
      """
    And file "conflicting_file" still contains unresolved conflicts

  Scenario: undo with unresolved conflict fails due to unresolved merge conflicts
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH   | COMMAND               |
      | new      | git add -A            |
      |          | git stash             |
      |          | git checkout existing |
      | existing | git branch -D new     |
      |          | git stash pop         |
    And it prints the error:
      """
      conflicts between your uncommmitted changes and the main branch
      """
    And the current branch is now "existing"

  Scenario: resolve and undo
    Given I resolve the conflict in "conflicting_file"
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH   | COMMAND               |
      | new      | git add -A            |
      |          | git stash             |
      |          | git checkout existing |
      | existing | git branch -D new     |
      |          | git stash pop         |
    And it prints the error:
      """
      conflicts between your uncommmitted changes and the main branch
      """
    And it does not print "to go back to where you started, run \"git-town undo\""
    And the current branch is now "existing"
    And the initial commits exist
    And file "conflicting_file" still has content "resolved content"

  Scenario: continue with unresolved conflict
    When I run "git-town continue"
    Then it prints the error:
      """
      you must resolve the conflicts before continuing
      """

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue" and close the editor
    Then it runs no commands
    And the current branch is now "new"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE            | FILE NAME        | FILE CONTENT |
      | main   | local, origin | conflicting commit | conflicting_file | main content |
      | new    | local         | conflicting commit | conflicting_file | main content |
    And file "conflicting_file" still has content "resolved content"

  Scenario: resolve, continue, and undo undoes the hack but cannot get back to the original branch due to merge conflicts
    Given I resolve the conflict in "conflicting_file"
    And I run "git-town continue" and close the editor
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH   | COMMAND               |
      | new      | git add -A            |
      |          | git stash             |
      |          | git checkout existing |
      | existing | git branch -D new     |
      |          | git stash pop         |
    And it prints the error:
      """
      conflicts between your uncommmitted changes and the main branch
      """
    And the current branch is now "existing"
    And the initial commits exist
    And the initial branches and lineage exist
    And file "conflicting_file" still has content "resolved content"
