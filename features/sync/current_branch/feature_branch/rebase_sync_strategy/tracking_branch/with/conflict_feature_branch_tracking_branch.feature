@skipWindows
Feature: handle conflicts between the current feature branch and its tracking branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | feature | local    | conflicting local commit  | conflicting_file | local content  |
      |         | origin   | conflicting origin commit | conflicting_file | origin content |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "feature"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                              |
      | feature | git fetch --prune --tags                             |
      |         | git push --force-with-lease --force-if-includes      |
      |         | git -c rebase.updateRefs=false rebase origin/feature |
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
    And no rebase is now in progress
    And the initial commits exist now

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
      | feature | GIT_EDITOR=true git rebase --continue           |
      |         | git push --force-with-lease --force-if-includes |
    And no rebase is now in progress
    And all branches are now synchronized
    And these committed files exist now
      | BRANCH  | NAME             | CONTENT          |
      | feature | conflicting_file | resolved content |

  Scenario: resolve, commit, and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git commit --no-edit"
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | GIT_EDITOR=true git rebase --continue           |
      |         | git push --force-with-lease --force-if-includes |
