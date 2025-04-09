@skipWindows
Feature: handle conflicts between the current feature branch and the main branch

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
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main --no-update-refs |
      |         | git push                                |
      |         | git checkout feature                    |
      | feature | git rebase main --no-update-refs        |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And a rebase is now in progress

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND            |
      | feature | git rebase --abort |
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
    And a rebase is now in progress

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue" and enter "resolved commit" for the commit message
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git -c core.editor=true rebase --continue       |
      |         | git push --force-with-lease --force-if-includes |
    And no rebase is now in progress
    And all branches are now synchronized
    And these committed files exist now
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | main content     |
      | feature | conflicting_file | resolved content |

  Scenario: resolve resulting in no changes and continue
    When I resolve the conflict in "conflicting_file" with "feature content"
    And I run "git-town continue" and enter "resolved commit" for the commit message
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git -c core.editor=true rebase --continue       |
      |         | git push --force-with-lease --force-if-includes |
    And all branches are now synchronized
    And no merge is in progress
    And these committed files exist now
      | BRANCH  | NAME             | CONTENT         |
      | main    | conflicting_file | main content    |
      | feature | conflicting_file | feature content |

  Scenario: resolve, commit, and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git rebase --continue" and enter "resolved commit" for the commit message
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git push --force-with-lease --force-if-includes |
    And all branches are now synchronized
    And no merge is in progress
    And these committed files exist now
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | main content     |
      | feature | conflicting_file | resolved content |
