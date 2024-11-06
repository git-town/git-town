@skipWindows
Feature: handle conflicts between the current feature branch and the main branch when using the "compress" strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
    And the current branch is "feature"
    And Git Town setting "sync-feature-strategy" is "compress"
    And an uncommitted file
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git add -A                              |
      |         | git stash                               |
      |         | git checkout main                       |
      | main    | git rebase origin/main --no-update-refs |
      |         | git push                                |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff main           |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And Git Town prints the error:
      """
      To continue after having resolved conflicts, run "git town continue".
      To go back to where you started, run "git town undo".
      To continue by skipping the current branch, run "git town skip".
      """
    And the current branch is still "feature"
    And the uncommitted file is stashed
    And a merge is now in progress

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND           |
      | feature | git merge --abort |
      |         | git stash pop     |
    And the current branch is still "feature"
    And the uncommitted file still exists
    And no merge is in progress
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local, origin | conflicting main commit    | conflicting_file | main content    |
      | feature | local         | conflicting feature commit | conflicting_file | feature content |

  Scenario: continue with unresolved conflict
    When I run "git-town continue"
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And the current branch is still "feature"
    And the uncommitted file is stashed
    And a merge is now in progress

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue" and enter "resolved commit" for the commit message
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                    |
      | feature | git commit --no-edit                       |
      |         | git merge --no-edit --ff origin/feature    |
      |         | git reset --soft main                      |
      |         | git commit -m "conflicting feature commit" |
      |         | git push --force-with-lease                |
      |         | git stash pop                              |
    And all branches are now synchronized
    And the current branch is still "feature"
    And no merge is in progress
    And the uncommitted file still exists
    And these committed files exist now
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | main content     |
      | feature | conflicting_file | resolved content |

  Scenario: resolve resulting in no changes and continue
    When I resolve the conflict in "conflicting_file" with "feature content"
    And I run "git-town continue" and enter "resolved commit" for the commit message
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                    |
      | feature | git commit --no-edit                       |
      |         | git merge --no-edit --ff origin/feature    |
      |         | git reset --soft main                      |
      |         | git commit -m "conflicting feature commit" |
      |         | git push --force-with-lease                |
      |         | git stash pop                              |
    And the current branch is still "feature"
    And all branches are now synchronized
    And no merge is in progress
    And the uncommitted file still exists
    And these committed files exist now
      | BRANCH  | NAME             | CONTENT         |
      | main    | conflicting_file | main content    |
      | feature | conflicting_file | feature content |

  Scenario: resolve, commit, and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git commit" and enter "resolved commit" for the commit message
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                    |
      | feature | git merge --no-edit --ff origin/feature    |
      |         | git reset --soft main                      |
      |         | git commit -m "conflicting feature commit" |
      |         | git push --force-with-lease                |
      |         | git stash pop                              |
    And the current branch is still "feature"
    And all branches are now synchronized
    And no merge is in progress
    And the uncommitted file still exists
    And these committed files exist now
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | main content     |
      | feature | conflicting_file | resolved content |
