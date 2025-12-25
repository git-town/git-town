Feature: handle conflicts between the current contribution branch and its tracking branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | LOCATIONS     |
      | contribution | contribution | local, origin |
    And the commits
      | BRANCH       | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | contribution | local    | conflicting local commit  | conflicting_file | local content  |
      |              | origin   | conflicting origin commit | conflicting_file | origin content |
    And the current branch is "contribution"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH       | COMMAND                                                   |
      | contribution | git fetch --prune --tags                                  |
      |              | git -c rebase.updateRefs=false rebase origin/contribution |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And a rebase is now in progress

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH       | COMMAND            |
      | contribution | git rebase --abort |
    And no rebase is now in progress
    And the initial branches and lineage exist now
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
    And I run "git-town continue" and close the editor
    Then Git Town runs the commands
      | BRANCH       | COMMAND                               |
      | contribution | GIT_EDITOR=true git rebase --continue |
      |              | git push                              |
    And no rebase is now in progress
    And these commits exist now
      | BRANCH       | LOCATION      | MESSAGE                   |
      | contribution | local, origin | conflicting origin commit |
      |              |               | conflicting local commit  |
    And these committed files exist now
      | BRANCH       | NAME             | CONTENT          |
      | contribution | conflicting_file | resolved content |

  Scenario: resolve, finish the rebase, and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git rebase --continue" and close the editor
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH       | COMMAND  |
      | contribution | git push |
    And no rebase is now in progress
    And these commits exist now
      | BRANCH       | LOCATION      | MESSAGE                   |
      | contribution | local, origin | conflicting origin commit |
      |              |               | conflicting local commit  |
    And these committed files exist now
      | BRANCH       | NAME             | CONTENT          |
      | contribution | conflicting_file | resolved content |
