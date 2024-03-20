Feature: handle conflicts between the current feature branch and its tracking branch

  Background:
    Given the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | feature | local    | conflicting local commit  | conflicting_file | local content  |
      |         | origin   | conflicting origin commit | conflicting_file | origin content |
    And an uncommitted file
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      |         | git add -A                         |
      |         | git stash                          |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
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
    And a merge is now in progress

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND           |
      | feature | git merge --abort |
      |         | git stash pop     |
    And the current branch is still "feature"
    And the uncommitted file still exists
    And no merge is in progress
    And the initial commits exist
    And the initial branches and lineage exist

  Scenario: continue with unresolved conflict
    When I run "git-town continue"
    Then it runs no commands
    And it prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And the current branch is still "feature"
    And the uncommitted file is stashed
    And a merge is now in progress

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git commit --no-edit     |
      |         | git merge --no-edit main |
      |         | git push                 |
      |         | git stash pop            |
    And all branches are now synchronized
    And the current branch is still "feature"
    And no merge is in progress
    And the uncommitted file still exists
    And these committed files exist now
      | BRANCH  | NAME             | CONTENT          |
      | feature | conflicting_file | resolved content |

  Scenario: resolve, commit, and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git commit --no-edit"
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git merge --no-edit main |
      |         | git push                 |
      |         | git stash pop            |
