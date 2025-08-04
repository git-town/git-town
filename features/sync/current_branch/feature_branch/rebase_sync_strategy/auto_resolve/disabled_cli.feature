Feature: don't auto-resolve phantom merge conflicts

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-1 | feature | main     | local, origin |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION | MESSAGE                     | FILE NAME        | FILE CONTENT     |
      | main     | local    | conflicting main commit     | conflicting_file | main content     |
      | branch-1 | local    | commit 1                    | other_file       | content          |
      | branch-2 | local    | conflicting branch-2 commit | conflicting_file | branch-2 content |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And origin deletes the "branch-1" branch
    And the current branch is "branch-2"
    When I run "git-town sync --auto-resolve=0"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                    |
      | branch-2 | git fetch --prune --tags                                   |
      |          | git checkout main                                          |
      | main     | git -c rebase.updateRefs=false rebase origin/main          |
      |          | git push                                                   |
      |          | git checkout branch-2                                      |
      | branch-2 | git pull                                                   |
      |          | git -c rebase.updateRefs=false rebase --onto main branch-1 |
      |          | GIT_EDITOR=true git rebase --continue                      |
      |          | git push --force-with-lease                                |
    # TODO: it should not run the rebase-continue and force-push at the end
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And Git Town prints something like:
      """
      could not apply .* conflicting branch-2 commit
      """
    And a rebase is now in progress

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

  @this
  Scenario: resolve, continue the rebase, and continue the sync
    When I resolve the conflict in "conflicting_file" with "branch-2 content"
    And I run "git rebase --continue" and enter "resolved commit" for the commit message
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                     |
      | branch-2 | git push --force-with-lease |
      |          | git branch -D branch-1      |
    And no rebase is now in progress
