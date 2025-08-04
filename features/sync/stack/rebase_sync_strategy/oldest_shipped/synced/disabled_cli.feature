Feature: disable auto-resolve phantom merge conflicts via CLI flag

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME        | FILE CONTENT               |
      | main   | local, origin | main commit | conflicting_file | line 1\n\nline 2\n\nline 3 |
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE                     | FILE NAME        | FILE CONTENT                                   |
      | branch-1 | local, origin | conflicting branch-1 commit | conflicting_file | line 1\n\nline 2 changed by branch-1\n\nline 3 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION | MESSAGE                     | FILE NAME        | FILE CONTENT                                   |
      | branch-2 | local    | conflicting branch-2 commit | conflicting_file | line 1\n\nline 2\n\nline 3 changed by branch-2 |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And origin ships the "branch-1" branch using the "squash-merge" ship-strategy
    And the current branch is "branch-2"
    When I run "git-town sync --auto-resolve=0"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                    |
      | branch-2 | git fetch --prune --tags                                   |
      |          | git checkout main                                          |
      | main     | git -c rebase.updateRefs=false rebase origin/main          |
      |          | git checkout branch-2                                      |
      | branch-2 | git pull                                                   |
      |          | git -c rebase.updateRefs=false rebase --onto main branch-1 |
      |          | git push --force-with-lease                                |
      |          | git branch -D branch-1                                     |
    # TODO: it should not run the rebase-continue and force-push at the end
    And no rebase is now in progress

  Scenario: undo
    When I run "git town undo"
    Then Git Town runs no commands
    And the initial branches and lineage exist now
# And no rebase is now in progress
# TODO: make this work

  Scenario: continue with unresolved conflicts
    When I run "git town continue"
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      you must resolve the conflicts before continuing
      """

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file" with "branch-2 content"
    And I run "git-town continue" and enter "resolved commit" for the commit message
    Then Git Town runs the commands
      | BRANCH   | COMMAND                     |
      | branch-2 | git push --force-with-lease |
    And Git Town prints the error:
      """
      You are not currently on a branch.
      """
# And no rebase is now in progress
# TODO: it should not print an error here but finish the sync

  Scenario: resolve, continue the rebase, and continue the sync
    When I resolve the conflict in "conflicting_file" with "branch-2 content"
    And I run "git rebase --continue" and enter "resolved commit" for the commit message
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                     |
      | branch-2 | git push --force-with-lease |
      |          | git branch -D branch-1      |
    And no rebase is now in progress
