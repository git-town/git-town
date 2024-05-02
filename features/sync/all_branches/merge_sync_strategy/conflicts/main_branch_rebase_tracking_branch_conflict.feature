Feature: handle rebase conflicts between main branch and its tracking branch

  Background:
    Given a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE            | FILE NAME        | FILE CONTENT    |
      | main    | local    | local main commit  | conflicting_file | local content   |
      |         | origin   | origin main commit | conflicting_file | origin content  |
      | feature | local    | feature commit     | feature_file     | feature content |
    And the current branch is "main"
    And an uncommitted file
    When I run "git-town sync --all"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git add -A               |
      |        | git stash                |
      |        | git rebase origin/main   |
    And it prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And it prints the error:
      """
      To continue after having resolved conflicts, run "git town continue".
      To go back to where you started, run "git town undo".
      """
    And the uncommitted file is stashed
    And a rebase is now in progress

  @this
  Scenario: undo
    # And inspect the repo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND            |
      | main   | git rebase --abort |
      |        | git stash pop      |
    And the current branch is now "main"
    And the uncommitted file still exists
    And the initial commits exist

  Scenario: continue with unresolved conflict
    When I run "git-town continue"
    Then it runs no commands
    And it prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And the uncommitted file is stashed
    And a rebase is now in progress

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue" and close the editor
    Then it runs the commands
      | BRANCH  | COMMAND                                 |
      | main    | git rebase --continue                   |
      |         | git push                                |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff origin/feature |
      |         | git merge --no-edit --ff main           |
      |         | git push                                |
      |         | git checkout main                       |
      | main    | git push --tags                         |
      |         | git stash pop                           |
    And all branches are now synchronized
    And the current branch is now "main"
    And the uncommitted file still exists
    And no rebase is in progress
    And these committed files exist now
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | resolved content |
      | feature | conflicting_file | resolved content |
      |         | feature_file     | feature content  |

  Scenario: resolve, finish the rebase, and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git rebase --continue" and close the editor
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND                                 |
      | main    | git push                                |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff origin/feature |
      |         | git merge --no-edit --ff main           |
      |         | git push                                |
      |         | git checkout main                       |
      | main    | git push --tags                         |
      |         | git stash pop                           |
