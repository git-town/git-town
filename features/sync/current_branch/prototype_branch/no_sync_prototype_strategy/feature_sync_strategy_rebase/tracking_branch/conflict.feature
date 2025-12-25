Feature: handle conflicts between the current prototype branch and its tracking branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | PARENT | LOCATIONS     |
      | prototype | prototype | main   | local, origin |
    And the commits
      | BRANCH    | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | prototype | local    | conflicting local commit  | conflicting_file | local content  |
      |           | origin   | conflicting origin commit | conflicting_file | origin content |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "prototype"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                                |
      | prototype | git fetch --prune --tags                               |
      |           | git push --force-with-lease --force-if-includes        |
      |           | git -c rebase.updateRefs=false rebase origin/prototype |
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
    And a rebase is now in progress

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH    | COMMAND            |
      | prototype | git rebase --abort |
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
      | BRANCH    | COMMAND                                         |
      | prototype | GIT_EDITOR=true git rebase --continue           |
      |           | git push --force-with-lease --force-if-includes |
    And no rebase is now in progress
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE                   |
      | prototype | local, origin | conflicting origin commit |
      |           |               | conflicting local commit  |
    And these committed files exist now
      | BRANCH    | NAME             | CONTENT          |
      | prototype | conflicting_file | resolved content |

  Scenario: resolve, finish the rebase, and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git rebase --continue" and close the editor
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                         |
      | prototype | git push --force-with-lease --force-if-includes |
    And no rebase is now in progress
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE                   |
      | prototype | local, origin | conflicting origin commit |
      |           |               | conflicting local commit  |
    And these committed files exist now
      | BRANCH    | NAME             | CONTENT          |
      | prototype | conflicting_file | resolved content |
