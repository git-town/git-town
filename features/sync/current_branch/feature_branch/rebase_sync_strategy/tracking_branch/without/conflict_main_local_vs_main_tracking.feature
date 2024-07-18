Feature: handle conflicts between the main branch and its tracking branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    Given Git Town setting "sync-feature-strategy" is "rebase"
    And the current branch is "feature"
    And the commits
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | main   | local    | conflicting local commit  | conflicting_file | local content  |
      |        | origin   | conflicting origin commit | conflicting_file | origin content |
    And an uncommitted file
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
      |         | git add -A               |
      |         | git stash                |
      |         | git checkout main        |
      | main    | git rebase origin/main   |
    And it prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And it prints the error:
      """
      To continue after having resolved conflicts, run "git town continue".
      To go back to where you started, run "git town undo".
      """
    And a rebase is now in progress
    And the uncommitted file is stashed

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND              |
      | main    | git rebase --abort   |
      |         | git checkout feature |
      | feature | git stash pop        |
    And the current branch is still "feature"
    And the uncommitted file still exists
    And no rebase is in progress
    And the initial commits exist

  @messyoutput
  Scenario: undo through another sync invocation
    When I run "git-town sync" and enter into the dialog:
      | DIALOG            | KEYS    |
      | choose what to do | 2 enter |
    Then it prints:
      """
      Handle unfinished command: undo
      """
    And it runs the commands
      | BRANCH  | COMMAND              |
      | main    | git rebase --abort   |
      |         | git checkout feature |
      | feature | git stash pop        |
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
    And a rebase is now in progress
    And the uncommitted file is stashed

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue" and close the editor
    Then it runs the commands
      | BRANCH  | COMMAND                                         |
      | main    | git rebase --continue                           |
      |         | git push                                        |
      |         | git checkout feature                            |
      | feature | git rebase main                                 |
      |         | git push --force-with-lease --force-if-includes |
      |         | git stash pop                                   |
    And all branches are now synchronized
    And the current branch is still "feature"
    And no rebase is in progress
    And the uncommitted file still exists
    And these committed files exist now
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | resolved content |
      | feature | conflicting_file | resolved content |

  Scenario: resolve, finish the rebase, and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git rebase --continue" and close the editor
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND                                         |
      | main    | git push                                        |
      |         | git checkout feature                            |
      | feature | git rebase main                                 |
      |         | git push --force-with-lease --force-if-includes |
      |         | git stash pop                                   |
    And all branches are now synchronized
    And the current branch is still "feature"
    And no rebase is in progress
    And the uncommitted file still exists
    And these committed files exist now
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | resolved content |
      | feature | conflicting_file | resolved content |
