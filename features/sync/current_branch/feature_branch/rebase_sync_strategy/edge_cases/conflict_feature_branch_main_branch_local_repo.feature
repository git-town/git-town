@skipWindows
Feature: handle conflicts between the current feature branch and the main branch (in a local repo)

  Background:
    Given Git Town setting "sync-feature-strategy" is "rebase"
    And my repo does not have an origin
    And the current branch is a local feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
    And an uncommitted file
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND         |
      | feature | git add -A      |
      |         | git stash       |
      |         | git rebase main |
    And it prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And it prints the error:
      """
      To continue after having resolved conflicts, run "git town continue".
      To go back to where you started, run "git town undo".
      To continue by skipping the current branch, run "git-town skip".
      """
    And the current branch is still "feature"
    And the uncommitted file is stashed
    And a rebase is now in progress

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND            |
      | feature | git rebase --abort |
      |         | git stash pop      |
    And the current branch is still "feature"
    And the uncommitted file still exists
    And no rebase is in progress
    And the initial commits exist

  Scenario: continue with unresolved conflict
    When I run "git-town continue"
    Then it runs no commands
    And it prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And the current branch is still "feature"
    And the uncommitted file is stashed
    And a rebase is now in progress

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git town continue".and enter "resolved commit" for the commit message
    Then it runs the commands
      | BRANCH  | COMMAND               |
      | feature | git rebase --continue |
      |         | git stash pop         |
    And all branches are now synchronized
    And the current branch is still "feature"
    And no rebase is in progress
    And the uncommitted file still exists
    And these committed files exist now
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | main content     |
      | feature | conflicting_file | resolved content |

  Scenario: resolve, commit, and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git rebase --continue" and enter "resolved commit" for the commit message
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND       |
      | feature | git stash pop |
