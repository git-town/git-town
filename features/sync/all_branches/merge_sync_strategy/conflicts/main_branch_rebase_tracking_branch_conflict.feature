Feature: handle rebase conflicts between main branch and its tracking branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION | MESSAGE            | FILE NAME        | FILE CONTENT    |
      | main    | local    | local main commit  | conflicting_file | local content   |
      |         | origin   | origin main commit | conflicting_file | origin content  |
      | feature | local    | feature commit     | feature_file     | feature content |
    And the current branch is "main"
    And an uncommitted file
    When I run "git-town sync --all"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                 |
      | main   | git fetch --prune --tags                |
      |        | git add -A                              |
      |        | git stash                               |
      |        | git rebase origin/main --no-update-refs |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And Git Town prints the error:
      """
      To continue after having resolved conflicts, run "git town continue".
      To go back to where you started, run "git town undo".
      """
    And the uncommitted file is stashed
    And a rebase is now in progress

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND            |
      | main   | git rebase --abort |
      |        | git stash pop      |
    And the current branch is now "main"
    And the uncommitted file still exists
    And the initial commits exist now

  Scenario: continue with unresolved conflict
    When I run "git-town continue"
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And the uncommitted file is stashed
    And a rebase is now in progress

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue" and close the editor
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                   |
      | main    | git -c core.editor=true rebase --continue |
      |         | git push                                  |
      |         | git checkout feature                      |
      | feature | git merge --no-edit --ff main             |
      |         | git merge --no-edit --ff origin/feature   |
      |         | git push                                  |
      |         | git checkout main                         |
      | main    | git push --tags                           |
      |         | git stash pop                             |
    And all branches are now synchronized
    And the current branch is now "main"
    And the uncommitted file still exists
    And no rebase is now in progress
    And these committed files exist now
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | resolved content |
      | feature | conflicting_file | resolved content |
      |         | feature_file     | feature content  |

  Scenario: resolve, finish the rebase, and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git rebase --continue" and close the editor
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | main    | git push                                |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff main           |
      |         | git merge --no-edit --ff origin/feature |
      |         | git push                                |
      |         | git checkout main                       |
      | main    | git push --tags                         |
      |         | git stash pop                           |
