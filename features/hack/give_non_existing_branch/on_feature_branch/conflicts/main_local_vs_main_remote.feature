Feature: conflicts between the main branch and its tracking branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | existing | feature | main   | local     |
    And the commits
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | main   | local    | conflicting local commit  | conflicting_file | local content  |
      |        | origin   | conflicting origin commit | conflicting_file | origin content |
    And the current branch is "existing"
    When I run "git-town hack new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                           |
      | existing | git fetch --prune --tags                          |
      |          | git checkout main                                 |
      | main     | git -c rebase.updateRefs=false rebase origin/main |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And a rebase is now in progress

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND               |
      | main   | git rebase --abort    |
      |        | git checkout existing |
    And no rebase is now in progress
    And the initial commits exist now

  Scenario: continue with unresolved conflict
    When I run "git-town continue"
    Then Git Town prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And a rebase is now in progress

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue" and close the editor
    Then Git Town runs the commands
      | BRANCH | COMMAND                               |
      | main   | GIT_EDITOR=true git rebase --continue |
      |        | git push                              |
      |        | git checkout -b new                   |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                   |
      | main   | local, origin | conflicting origin commit |
      |        |               | conflicting local commit  |
    And these committed files exist now
      | BRANCH | NAME             | CONTENT          |
      | main   | conflicting_file | resolved content |
      | new    | conflicting_file | resolved content |

  Scenario: resolve, finish the rebase, and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git rebase --continue" and close the editor
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH | COMMAND             |
      | main   | git push            |
      |        | git checkout -b new |
